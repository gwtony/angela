package main

import (
	"fmt"
	"time"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/angela/handler"
)

func main() {
	err := api.Init("angela.conf")
	if err != nil {
		fmt.Println("[Error] Init api failed")
		return
	}
	config := api.GetConfig()
	log := api.GetLog()

	err = handler.InitContext(config, log)
	if err != nil {
		fmt.Println("[Error] Init angela failed")
		time.Sleep(time.Second)
		return
	}

	api.Run()
}
