package db

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocql/gocql"
)

type Author struct {
	ID        gocql.UUID `json:"id" db:"id" cql:"id"`
	Firstname string     `json:"firstName" db:"first_name" cql:"first_name"`
	Lastname  string     `json:"lastName" db:"last_name" cql:"last_name"`
}

type Message struct {
	Author    Author     `json:"author" db:"author" cql:"author"`
	ID        gocql.UUID `json:"id" db:"id" cql:"id"`
	CreatedAt int64      `json:"createdAt" db:"created_at" cql:"created_at"`
	Name      string     `json:"name" db:"name" cql:"name"`
	Size      int64      `json:"size" db:"size" cql:"size"`
	Status    string     `json:"status" db:"status" cql:"status"`
	Type      string     `json:"type" db:"type" cql:"type"`
	Uri       string     `json:"uri" db:"uri" cql:"uri"`
	Width     int64      `json:"width" db:"width" cql:"width"`
	Height    int64      `json:"height" db:"height" cql:"height"`
	Text      string     `json:"text" db:"text" cql:"text"`
	RoomId    gocql.UUID `json:"roomId" db:"room_id" cql:"room_id"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Restful API using Go and Cassandra!")
}

func CreateMessage(NewMessage Message) {

	if err := Session.Query("INSERT INTO message(author, id, created_at, name, size, status, type, uri, width, height, text, room_id) VALUES("+
		"{ id: ?, first_name: ?, last_name: ?}, "+
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		NewMessage.Author.ID, NewMessage.Author.Firstname, NewMessage.Author.Lastname,
		NewMessage.ID, NewMessage.CreatedAt, NewMessage.Name, NewMessage.Size, NewMessage.Status, NewMessage.Type, NewMessage.Uri,
		NewMessage.Width, NewMessage.Height, NewMessage.Text, NewMessage.RoomId).Exec(); err != nil {
		fmt.Println("Error while inserting")
		fmt.Println(err)
	}

}

func GetAllMessagesByRoomId(roomId gocql.UUID) []Message {
	var messages []Message

	scanner := Session.Query(
		"SELECT author, id, created_at, name, size, status, type, uri, width, height, text, room_id FROM message where room_id = ?",
		roomId,
	).Iter().Scanner()

	for scanner.Next() {
		var message Message
		err := scanner.Scan(&message.Author, &message.ID, &message.CreatedAt, &message.Name, &message.Size, &message.Status, &message.Type, &message.Uri,
			&message.Width, &message.Height, &message.Text, &message.RoomId)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return messages
}

func GetBotUser() Author {
	if os.Getenv("BOT_USER_UUID") == "" {
		os.Setenv("BOT_USER_UUID", "88a2c050-d1e8-11ee-bfde-16e722bb8d94")
	}
	uuid, err := gocql.ParseUUID(os.Getenv("BOT_USER_UUID"))
	if err != nil {
		log.Fatal("Cannot get uuid for bot user")
	}
	return Author{
		ID:        uuid,
		Firstname: "Bot",
		Lastname:  "Zanger",
	}
}

func InsertAndGetBotMessage(text string, roomId gocql.UUID) Message {
	message := Message{
		Author:    GetBotUser(),
		ID:        gocql.TimeUUID(),
		CreatedAt: time.Now().UnixMilli(),
		Name:      "",
		Size:      0,
		Status:    "seen",
		Type:      "text",
		Uri:       "",
		Width:     0,
		Height:    0,
		Text:      text,
		RoomId:    roomId,
	}
	CreateMessage(message)
	return message
}

// func GetOneStudent(w http.ResponseWriter, r *http.Request) {
// 	StudentID := mux.Vars(r)["id"]
// 	var students []Student
// 	m := map[string]interface{}{}

// 	iter := Session.Query("SELECT * FROM students WHERE id=?", StudentID).Iter()
// 	for iter.MapScan(m) {
// 		students = append(students, Student{
// 			ID:        m["id"].(int),
// 			Firstname: m["firstname"].(string),
// 			Lastname:  m["lastname"].(string),
// 			Age:       m["age"].(int),
// 		})
// 		m = map[string]interface{}{}
// 	}

// 	Conv, _ := json.MarshalIndent(students, "", " ")
// 	fmt.Fprintf(w, "%s", string(Conv))

// }

// func CountAllStudents(w http.ResponseWriter, r *http.Request) {

// 	var Count string
// 	err := Session.Query("SELECT count(*) FROM students").Scan(&Count)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Fprintf(w, "%s ", Count)

// }
// func DeleteOneStudent(w http.ResponseWriter, r *http.Request) {
// 	StudentID := mux.Vars(r)["id"]
// 	if err := Session.Query("DELETE FROM students WHERE id = ?", StudentID).Exec(); err != nil {
// 		fmt.Println("Error while deleting")
// 		fmt.Println(err)
// 	}
// 	fmt.Fprintf(w, "deleted successfully the student num %s ", StudentID)
// }

// func DeleteAllStudents(w http.ResponseWriter, r *http.Request) {

// 	if err := Session.Query("TRUNCATE students").Exec(); err != nil {
// 		fmt.Println("Error while deleting all students")
// 		fmt.Println(err)
// 	}
// 	fmt.Fprintf(w, "deleted all successfully")

// }

// func UpdateStudent(w http.ResponseWriter, r *http.Request) {
// 	StudentID := mux.Vars(r)["id"]
// 	var UpdateStudent Student
// 	reqBody, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data properly")
// 	}
// 	json.Unmarshal(reqBody, &UpdateStudent)
// 	if err := Session.Query("UPDATE students SET firstname = ?, lastname = ?, age = ? WHERE id = ?",
// 		UpdateStudent.Firstname, UpdateStudent.Lastname, UpdateStudent.Age, StudentID).Exec(); err != nil {
// 		fmt.Println("Error while updating")
// 		fmt.Println(err)
// 	}
// 	fmt.Fprintf(w, "updated successfully")

// }
