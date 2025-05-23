package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types/events"
	api "github.com/faisolarifin/wacoregateway/http/grpc"
	"github.com/faisolarifin/wacoregateway/provider"
	"github.com/faisolarifin/wacoregateway/util"
	"github.com/go-playground/validator/v10"
	"go.mau.fi/whatsmeow"
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
	validate := validator.New()
	logger.Infof(provider.AppLog, "Application started")

	// ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, "MAIN")

	go func(logger provider.ILogger) {
		app := api.NewApp(validate, logger)

		addr := fmt.Sprintf(":%v", util.Configuration.Server.Port)
		server, err := app.GRPCServer()
		if err != nil {
			log.Fatal(err)
		}

		lis, err := net.Listen("tcp", ":"+addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		logger.Infof(provider.AppLog, "gRPC server listening on :%v", addr)
		if err := server.Serve(lis); err != nil {
			logger.Errorf(provider.AppLog, "failed to serve: %v", err)
		}
	}(logger)

	// sqlDB, err := provider.NewPostgresConnection()
	// if err != nil {
	// 	logger.Errorf(provider.AppLog, "Failed to connect to database:", err)
	// }
	// defer sqlDB.Close()

	// dbLog := waLog.Stdout("Database", "DEBUG", true)
	// container := sqlstore.NewWithDB(sqlDB, "postgres", dbLog)

	// clientLog := waLog.Stdout("Client", "DEBUG", true)

	// devices, err := container.GetAllDevices(ctx)
	// if err != nil {
	// 	logger.Errorf(provider.AppLog, "Failed to get device store: %v", err)
	// }

	// for _, dev := range devices {
	// 	client := whatsmeow.NewClient(dev, clientLog)
	// 	err := client.Connect()
	// 	if err != nil {
	// 		fmt.Printf("failed to connect device %s: %v\n", dev.ID.String(), err)
	// 		continue
	// 	}
	// 	clients[dev.ID.String()] = client
	// }

	// r := gin.Default()

	// r.GET("/devices", func(c *gin.Context) {
	// 	devices := []string{}
	// 	for id := range clients {
	// 		devices = append(devices, id)
	// 	}
	// 	c.JSON(200, gin.H{"devices": devices})
	// })

	// r.POST("/devices/new/:number", func(c *gin.Context) {
	// 	num := c.Param("number")
	// 	jid := types.NewJID(num, types.DefaultUserServer)
	// 	device := container.NewDevice()

	// 	client := whatsmeow.NewClient(device, clientLog)
	// 	clients[jid.String()] = client

	// 	if client.Store.ID == nil {
	// 		qrChan, _ := client.GetQRChannel(context.Background())
	// 		go func() {
	// 			_ = client.Connect()
	// 		}()

	// 		for evt := range qrChan {
	// 			if evt.Event == "code" {
	// 				c.JSON(200, gin.H{"qr": evt.Code})
	// 				return
	// 			}
	// 		}
	// 	} else {
	// 		c.JSON(200, gin.H{"message": "Device already logged in."})
	// 	}
	// })

	// r.POST("/send", func(c *gin.Context) {
	// 	type SendMessageRequest struct {
	// 		SenderJID   string `json:"sender_jid"`
	// 		To          string `json:"to"`
	// 		MessageText string `json:"message"`
	// 	}

	// 	var req SendMessageRequest
	// 	if err := c.ShouldBindJSON(&req); err != nil {
	// 		c.JSON(400, gin.H{"error": "invalid request"})
	// 		return
	// 	}

	// 	client, exists := clients[req.SenderJID]
	// 	if !exists {
	// 		c.JSON(404, gin.H{"error": "sender device not found"})
	// 		return
	// 	}

	// 	jid, err := types.ParseJID(req.To)
	// 	if err != nil {
	// 		c.JSON(400, gin.H{"error": "invalid recipient JID"})
	// 		return
	// 	}

	// 	resp, err := client.SendMessage(context.Background(), jid, &waE2E.Message{
	// 		Conversation: proto.String(req.MessageText),
	// 	})

	// 	if err != nil {
	// 		c.JSON(500, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(200, gin.H{"status": "message sent " + resp.ID})
	// })

	// r.GET("/contacts/:sender_jid", func(c *gin.Context) {
	// 	senderJID := c.Param("sender_jid")
	// 	client, exists := clients[senderJID]
	// 	if !exists {
	// 		c.JSON(404, gin.H{"error": "client not found"})
	// 		return
	// 	}

	// 	contacts, err := client.Store.Contacts.GetAllContacts(ctx)
	// 	if err != nil {
	// 		c.JSON(500, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	result := []map[string]string{}
	// 	for jid, contact := range contacts {
	// 		result = append(result, map[string]string{
	// 			"jid":   jid.String(),
	// 			"name":  contact.FirstName,
	// 			"short": contact.FullName,
	// 		})
	// 	}
	// 	c.JSON(200, result)
	// })

	// r.GET("/groups/:sender_jid", func(c *gin.Context) {
	// 	senderJID := c.Param("sender_jid")
	// 	client, exists := clients[senderJID]
	// 	if !exists {
	// 		c.JSON(404, gin.H{"error": "client not found"})
	// 		return
	// 	}

	// 	groups, err := client.GetJoinedGroups()
	// 	if err != nil {
	// 		c.JSON(500, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	result := []map[string]string{}
	// 	for _, group := range groups {
	// 		result = append(result, map[string]string{
	// 			"jid":  group.JID.String(),
	// 			"name": group.Name,
	// 		})
	// 	}
	// 	c.JSON(200, result)
	// })

	// r.Run(":8080")

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infof(provider.AppLog, "Receiving signal: %s", sig)
}

func protoString(s string) *string {
	return &s
}
