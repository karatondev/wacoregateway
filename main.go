package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types/events"
	"github.com/faisolarifin/wacoregateway/model/constant"
	"github.com/faisolarifin/wacoregateway/provider"
	"github.com/faisolarifin/wacoregateway/util"
	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var clients = make(map[string]*whatsmeow.Client)

func init() {
	if err := util.LoadConfig("."); err != nil {
		log.Fatal(err)
	}
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Handle incoming messages
		fmt.Println(v)
		// sender := v.Info.Sender.String()
		// content := v.Message.GetConversation()
		// fmt.Printf("New message from %s: %s\n", sender, content)
	}
}

func main() {
	logger := provider.NewLogger()
	logger.Infof(provider.AppLog, "Application started")

	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, "MAIN")

	// Create database connection
	sqlDB, err := provider.NewPostgresConnection()
	if err != nil {
		logger.Errorf(provider.AppLog, "Failed to connect to database:", err)
	}
	defer sqlDB.Close()

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container := sqlstore.NewWithDB(sqlDB, "postgres", dbLog)

	// deviceStore, err := container.GetFirstDevice(ctx)
	// if err != nil {
	// 	logger.Errorf("Failed to get device store: %v", err)
	// }

	clientLog := waLog.Stdout("Client", "DEBUG", true)

	devices, err := container.GetAllDevices(ctx)
	if err != nil {
		logger.Errorf(provider.AppLog, "Failed to get device store: %v", err)
	}

	for _, dev := range devices {
		client := whatsmeow.NewClient(dev, clientLog)
		err := client.Connect()
		if err != nil {
			fmt.Printf("failed to connect device %s: %v\n", dev.ID.String(), err)
			continue
		}
		clients[dev.ID.String()] = client
	}

	r := gin.Default()

	r.GET("/devices", func(c *gin.Context) {
		devices := []string{}
		for id := range clients {
			devices = append(devices, id)
		}
		c.JSON(200, gin.H{"devices": devices})
	})

	r.POST("/devices/new/:number", func(c *gin.Context) {
		num := c.Param("number")
		jid := types.NewJID(num, types.DefaultUserServer)
		device := container.NewDevice()

		client := whatsmeow.NewClient(device, clientLog)
		clients[jid.String()] = client

		if client.Store.ID == nil {
			qrChan, _ := client.GetQRChannel(context.Background())
			go func() {
				_ = client.Connect()
			}()

			for evt := range qrChan {
				if evt.Event == "code" {
					c.JSON(200, gin.H{"qr": evt.Code})
					return
				}
			}
		} else {
			c.JSON(200, gin.H{"message": "Device already logged in."})
		}
	})

	r.POST("/send", func(c *gin.Context) {
		type SendMessageRequest struct {
			SenderJID   string `json:"sender_jid"`
			To          string `json:"to"`
			MessageText string `json:"message"`
		}

		var req SendMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		client, exists := clients[req.SenderJID]
		if !exists {
			c.JSON(404, gin.H{"error": "sender device not found"})
			return
		}

		jid, err := types.ParseJID(req.To)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid recipient JID"})
			return
		}

		resp, err := client.SendMessage(context.Background(), jid, &waE2E.Message{
			Conversation: proto.String(req.MessageText),
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// jid := types.NewJID(req.To, types.DefaultUserServer)
		// msg := &waProto.Message{
		// 	Conversation: protoString(req.MessageText),
		// }

		// _, err := client.SendMessage(context.Background(), jid, "", msg)
		// if err != nil {
		// 	c.JSON(500, gin.H{"error": err.Error()})
		// 	return
		// }

		c.JSON(200, gin.H{"status": "message sent " + resp.ID})
	})

	r.Run(":8080")

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infof(provider.AppLog, "Receiving signal: %s", sig)
}

func protoString(s string) *string {
	return &s
}
