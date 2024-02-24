package ws

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"

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
	RoomId    gocql.UUID
}

type Server struct {
	clients       map[*websocket.Conn]Client
	handleMessage func(bytes []byte) // хандлер новых сообщений
}

func StartServer(handleMessage func(bytes []byte)) *Server {
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

	header := r.Header.Get("Sec-WebSocket-Protocol")
	if header != "" {
		uuid, err := gocql.ParseUUID(header)
		if err != nil {
			fmt.Println("An error occurred while getting header room id")
			fmt.Println(err)
		}

		server.clients[connection] = Client{
			RoomId:    uuid,
			IsEnabled: true,
		} // Сохраняем соединение, используя его как ключ
		defer delete(server.clients, connection) // Удаляем соединение

		messages := db.GetAllMessagesByRoomId(uuid)

		if len(messages) != 0 {
			for _, message := range messages {
				server.WriteMessage(message)
			}
		} else {
			server.WriteMessage(db.InsertAndGetBotMessage("Hello, there is lawyer bot", uuid))
		}
		for {
			mt, message, err := connection.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go server.handleMessage(message)
		}
	}
}

func (server *Server) WriteMessage(message db.Message) {
	Conv, _ := json.MarshalIndent(message, "", " ")
	fmt.Println(string(Conv))
	for conn, client := range server.clients {
		if client.RoomId == message.RoomId {
			conn.WriteMessage(websocket.TextMessage, Conv)
			db.CreateMessage(message)
		}
	}
}
