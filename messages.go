package main

import (
	"db"
	"encoding/json"
	"fmt"
	"time"
	"ws"
)

func main() {
	ws.StartServer(messageHandler)

	for {
		time.Sleep(time.Second)
	}
}

func messageHandler(bytes []byte) {
	var NewMessage db.Message
	err := json.Unmarshal(bytes, &NewMessage)
	if err != nil {
		fmt.Println("An error occured while unmarshalling message ", string(bytes))
	}

	if NewMessage.Type != "" {
		db.CreateMessage(NewMessage)
		Conv, _ := json.MarshalIndent(NewMessage, "", " ")
		fmt.Println(string(Conv))
	}
}
