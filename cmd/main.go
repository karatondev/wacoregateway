package main

import (
	"log"
	"wacoregateway/internal/app"
	"wacoregateway/util"
)

func main() {
	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}
