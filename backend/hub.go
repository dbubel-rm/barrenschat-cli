package main

import (
	"github.com/engineerbeard/barrenschat/httpscerts"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	b "github.com/engineerbeard/barrenschat/shared"
	"strings"
	"fmt"
)

const CERT_PEM = "cert.pem"
const KEY_PEM = "key.pem"
const HUB_ADDR = "localhost:8081"


var upgrader = websocket.Upgrader{EnableCompression:true} // use default options
var connectedClients map[string][]BChatClient

type BChatClient struct {
	Name   string
	Room   string
	WsConn *websocket.Conn
	Uid    string
}

func init() {
	rand.Seed(time.Now().UnixNano())
	connectedClients = make(map[string][]BChatClient)
	connectedClients[b.MAIN_ROOM] = []BChatClient{}
}

func BroadcastMessage(room string, s b.BMessage) {
	log.Println("Broadcasting msg:", strings.Replace(fmt.Sprint(s), "\n", "", -1))
	for _, j := range connectedClients[room] {
		go j.WsConn.WriteJSON(s)
	}
}

func FindClient(id string) (string, int, string) {
	var room string
	var name string
	var idx int
	for _, j := range connectedClients {
		for x := 0; x < len(j); x++ {
			if strings.Contains(j[x].Uid, id) {
				return j[x].Room, x, j[x].Name
			}
		}
	}
	return room, idx, name
}

func GetNamesInRoom(room string) (string) {
	var nameList string
	for _, j := range connectedClients[room] {
		nameList = fmt.Sprint(nameList, j.Name, "\n")
	}
	return nameList
}

func handleNewClient(c *websocket.Conn) {
	defer c.Close()
	var id string
	var bMessage b.BMessage

	c.ReadJSON(&bMessage)
	//c.WriteJSON(b.BMessage{
	//	MsgType: b.B_CONNECT,
	//	Name: bMessage.Payload,
	//	Room: b.MAIN_ROOM,
	//	Payload: "Connected to Barrenschat ",
	//	TimeStamp: time.Now(),
	//	Data:GetNamesInRoom(b.MAIN_ROOM),
	//	})
	BClient := BChatClient{WsConn: c, Name: bMessage.Payload, Room: b.MAIN_ROOM, Uid: bMessage.Uid}
	connectedClients[b.MAIN_ROOM] = append(connectedClients[b.MAIN_ROOM], BClient)
	c.WriteJSON(b.BMessage{
		MsgType:   b.B_CONNECT,
		Name:      bMessage.Payload,
		Room:      b.MAIN_ROOM,
		Payload:   "Welcome to BarrensChat\nCommands:\n /name <name>\n /room <room>\n /newroom <newroom>",
		TimeStamp: time.Now(),
		Data:      GetNamesInRoom(b.MAIN_ROOM),
	})
	id = bMessage.Uid
	BroadcastMessage(b.MAIN_ROOM, b.BMessage{
		MsgType:   b.B_CONNECT,
		Name:      bMessage.Payload,
		Room:      b.MAIN_ROOM,
		Payload:   "New Connection!",
		TimeStamp: time.Now(),
		Data:      GetNamesInRoom(b.MAIN_ROOM),
	})



	for {
		err := c.ReadJSON(&bMessage)
		if err == nil { // Process message

			room, idx, _ := FindClient(id)

			if bMessage.MsgType == b.B_NAMECHANGE {
				connectedClients[room][idx].Name = bMessage.Name
				bMessage.Data = GetNamesInRoom(room)
			}
			BroadcastMessage(room, bMessage)

		} else { // Clean up
			BClient.WsConn.Close()
			room, idx, name := FindClient(id)
			connectedClients[room] = append(connectedClients[room][:idx], connectedClients[room][idx+1:]...)
			BroadcastMessage(room, b.BMessage{
				MsgType:   b.B_DISCONNECT,
				TimeStamp: time.Now(),
				Payload:   name + " Disconnected",
				Data:      GetNamesInRoom(room),
				})
			log.Println(err) // Connection is over
			break
		}
	}
}

func WsStart(s string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		c, _ := upgrader.Upgrade(w, r, nil)
		go handleNewClient(c)
	}
}

func main() {
	f, err := os.OpenFile("hub_log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stdout)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err = httpscerts.Generate(CERT_PEM, KEY_PEM, HUB_ADDR)
	if err != nil {
		log.Fatal("Error: Couldn't create https certs.")
	}

	router := httprouter.New()
	router.GET("/bchatws", WsStart("s"))
	log.Println("Server started on:", HUB_ADDR)
	log.Println(http.ListenAndServeTLS(HUB_ADDR, CERT_PEM, KEY_PEM, router))

}
