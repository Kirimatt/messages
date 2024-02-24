package main

import (
	"bytes"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"ws"
)

type ModelResponse struct {
	OutputText string `json:"output_text"`
}

var server *ws.Server

func init() {
	if os.Getenv("MODEL_URL") == "" {
		os.Setenv("MODEL_URL", "http://20.163.183.142")
	}
	if os.Getenv("MODEL_CHAT_API") == "" {
		os.Setenv("MODEL_CHAT_API", "/chat")
	}
	fmt.Println("Values setted successfuly")
}

func main() {
	server = ws.StartServer(messageHandler)

	for {
		time.Sleep(time.Second)
	}
}

func messageHandler(messageBytes []byte) {
	var NewMessage db.Message
	err := json.Unmarshal(messageBytes, &NewMessage)
	if err != nil {
		fmt.Println("An error occured while unmarshalling message ", string(messageBytes))
	}

	if NewMessage.Type != "" {
		db.CreateMessage(NewMessage)
		Conv, _ := json.MarshalIndent(NewMessage, "", " ")
		fmt.Println(string(Conv))

		server.WriteMessage(
			db.InsertAndGetBotMessage(
				postToModel(NewMessage).OutputText,
				NewMessage.RoomId,
			),
		)
	}
}

func postToModel(NewMessage db.Message) ModelResponse {
	posturl := os.Getenv("MODEL_URL") + os.Getenv("MODEL_CHAT_API")

	body := []byte(
		fmt.Sprintf(
			`{"input_text": "%s"}`,
			NewMessage.Text,
		),
	)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	ModelResponse := &ModelResponse{}
	derr := json.NewDecoder(res.Body).Decode(ModelResponse)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}
	return *ModelResponse
}
