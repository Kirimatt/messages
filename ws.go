package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type Client struct {
	IsEnabled bool
	RoomId    *gocql.UUID
}

type Server struct {
	clients       map[*websocket.Conn]Client
	handleMessage func(bytes []byte) // хандлер новых сообщений
}

type AuthorizedUser struct {
	Authorization string `json:"authorization"`
}

func StartServer(handleMessage func(bytes []byte)) *Server {
	if os.Getenv("BOT_FIRST_MESSAGE") == "" {
		os.Setenv("BOT_FIRST_MESSAGE", "Hello, there is lawyer bot")
	}

	server := Server{
		make(map[*websocket.Conn]Client),
		handleMessage,
	}

	http.HandleFunc("/", server.echo)
	go http.ListenAndServe("localhost:8080", nil) // Уводим http сервер в горутину

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	server.clients[connection] = Client{
		RoomId:    nil,
		IsEnabled: true,
	} // Сохраняем соединение, используя его как ключ, выставляем номер комнаты nil
	defer delete(server.clients, connection) // Удаляем соединение

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		var AuthorizedUser AuthorizedUser
		json.Unmarshal(message, &AuthorizedUser)
		if AuthorizedUser.Authorization != "" {
			server.authorize(AuthorizedUser, connection) // выставляем номер комнаты uuid
		} else {
			go server.handleMessage(message)
		}
	}
}

func (server *Server) authorize(AuthorizedUser AuthorizedUser, connection *websocket.Conn) {
	uuid, err := gocql.ParseUUID(AuthorizedUser.Authorization)
	if err != nil {
		fmt.Println("An error occurred while getting authorization uuid")
		fmt.Println(err)
	}

	server.clients[connection] = Client{
		RoomId:    &uuid,
		IsEnabled: true,
	}

	messages := GetAllMessagesByRoomId(uuid)

	if len(messages) != 0 {
		for _, message := range messages {
			server.WriteMessage(message)
		}
	} else {
		server.WriteMessage(InsertAndGetBotMessage(os.Getenv("BOT_FIRST_MESSAGE"), uuid))
	}
}

func (server *Server) WriteMessage(message Message) {
	Conv, _ := json.MarshalIndent(message, "", " ")
	fmt.Println(string(Conv))
	for conn, client := range server.clients {
		if *client.RoomId == message.RoomId {
			conn.WriteMessage(websocket.TextMessage, Conv)
			CreateMessage(message)
		}
	}
}
