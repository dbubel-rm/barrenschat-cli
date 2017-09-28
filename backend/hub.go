package main

import (
	"github.com/engineerbeard/barrenschat/httpscerts"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"fmt"
	b "github.com/engineerbeard/barrenschat/shared"
	"strings"
	"sync"
)

const CERT_PEM = "cert.pem"
const KEY_PEM = "key.pem"
const HUB_ADDR = "localhost:8081"

//var connectedClients map[string][]BChatClient
var connectedClients ServerStruct

type ServerStruct struct {
	Clients map[string][]b.BChatClient
	mu      sync.Mutex
}

func (s *ServerStruct) ChangeClientName(room string, index int, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[room][index].ChangeName(name)
}
func (s *ServerStruct) ChangeClientRoom(room string, idx int, bMessage b.BMessage) {

	s.Clients[room][idx].ChangeRoom(bMessage.Room)
	s.mu.Lock()

	if _, room_exists := s.Clients[bMessage.Room]; room_exists {
		s.Clients[bMessage.Room] = append(s.Clients[bMessage.Room], s.Clients[room][idx])
	} else {
		s.Clients[bMessage.Room] = []b.BChatClient{s.Clients[room][idx]}
	}

	s.Clients[room] = append(s.Clients[room][:idx], s.Clients[room][idx+1:]...)
	s.mu.Unlock()

	// Update clients
	bMessage.RoomData = s.GetRoomsString()
	bMessage.OnlineData = s.GetNamesInRoom(bMessage.Room)
	bMessage.Payload = fmt.Sprint(bMessage.Name, " joined the room.")
	s.BroadcastMessage(bMessage.Room, bMessage) // Broadcast to new room
	bMessage.Payload = fmt.Sprint(bMessage.Name, " left the room.")
	bMessage.OnlineData = s.GetNamesInRoom(room)
	s.BroadcastMessage(room, bMessage)
}

func (s *ServerStruct) FindClient(id string) (string, int, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var room string
	var name string
	var idx int
	for _, j := range s.Clients {
		for x := 0; x < len(j); x++ {
			if strings.Contains(j[x].Uid, id) {
				return j[x].Room, x, j[x].Name
			}
		}
	}
	return room, idx, name
}

func (s *ServerStruct) GetNamesInRoom(room string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var nameList string
	for i := range s.Clients[room] {
		nameList = fmt.Sprint(nameList, s.Clients[room][i].Name, "\n")
	}
	return nameList
}

func (s *ServerStruct) GetRoomsString() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var keys string
	for k := range s.Clients {
		keys = keys + k + "\n"
	}
	return keys
}

func (s *ServerStruct) AddRoom(r string) {
	s.mu.Lock()
	s.Clients[b.MAIN_ROOM] = []b.BChatClient{}
	s.mu.Unlock()
}
func (s *ServerStruct) BroadcastMessage(r string, m b.BMessage) {
	s.mu.Lock()
	for i := range s.Clients[r] {
		s.Clients[r][i].SendMessage(m)
	}
	s.mu.Unlock()
}

func (s *ServerStruct) AddClient(c *b.BChatClient) {
	s.mu.Lock()
	s.Clients[b.MAIN_ROOM] = append(s.Clients[b.MAIN_ROOM], *c)
	s.mu.Unlock()
}

func (s *ServerStruct) RemoveClient(room string, idx int) {
	s.mu.Lock()
	s.Clients[room] = append(s.Clients[room][:idx], s.Clients[room][idx+1:]...)
	s.mu.Unlock()

}

func init() {
	rand.Seed(time.Now().UnixNano())
	connectedClients.Clients = make(map[string][]b.BChatClient)
	connectedClients.mu = sync.Mutex{}
	connectedClients.AddRoom(b.MAIN_ROOM)
	//
}

//func BroadcastMessage(room string, s b.BMessage) {
//	//log.Println(room)
//	//log.Println(strings.Replace(fmt.Sprint(s), "\n", "", -1))
//
//}

func handleNewClient(c *websocket.Conn) {
	defer c.Close()
	var id string
	//var bMessage b.BMessage

	//c.ReadJSON(&bMessage)
	BClient := b.BChatClient{WsConn: c}
	bMessage, _ := BClient.ReadMessage()
	//Name: bMessage.Payload, Room: b.MAIN_ROOM, Uid: bMessage.Uid

	BClient.Name = bMessage.Payload
	BClient.Room = b.MAIN_ROOM
	BClient.Uid = bMessage.Uid

	connectedClients.AddClient(&BClient)
	BClient.SendMessage(b.BMessage{
		MsgType:    b.B_CONNECT,
		Name:       bMessage.Payload,
		Room:       b.MAIN_ROOM,
		Payload:    "Welcome to BarrensChat\nCommands:\n /name <name>\n /room <room>\n",
		TimeStamp:  time.Now(),
		OnlineData: connectedClients.GetNamesInRoom(b.MAIN_ROOM),
		RoomData:   connectedClients.GetRoomsString(),
	})
	id = bMessage.Uid
	connectedClients.BroadcastMessage(b.MAIN_ROOM, b.BMessage{
		MsgType:    b.B_CONNECT,
		Name:       bMessage.Payload,
		Room:       b.MAIN_ROOM,
		Payload:    "New Connection!",
		TimeStamp:  time.Now(),
		OnlineData: connectedClients.GetNamesInRoom(b.MAIN_ROOM),
		RoomData:   connectedClients.GetRoomsString(),
	})

	for {
		//err := c.ReadJSON(&bMessage)
		bMessage, err := BClient.ReadMessage()
		if err == nil { // Process message

			room, idx, _ := connectedClients.FindClient(id)
			bMessage.OnlineData = connectedClients.GetNamesInRoom(room)
			if bMessage.MsgType == b.B_MESSAGE {
				connectedClients.BroadcastMessage(room, bMessage)
			} else if bMessage.MsgType == b.B_NAMECHANGE {
				connectedClients.ChangeClientName(room, idx, bMessage.Name)
				bMessage.OnlineData = connectedClients.GetNamesInRoom(room)
				connectedClients.BroadcastMessage(room, bMessage)
			} else if bMessage.MsgType == b.B_ROOMCHANGE {
				connectedClients.ChangeClientRoom(room, idx,bMessage)

			}
		} else { // Clean up
			BClient.Close()
			room, idx, name := connectedClients.FindClient(id)
			connectedClients.RemoveClient(room, idx)
			connectedClients.BroadcastMessage(room, b.BMessage{
				MsgType:    b.B_DISCONNECT,
				TimeStamp:  time.Now(),
				Payload:    name + " Disconnected",
				OnlineData: connectedClients.GetNamesInRoom(room),
				RoomData:   connectedClients.GetRoomsString(),
			})
			//log.Println(err) // Connection is over
			break
		}
	}
}

func WsStart(upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := upgrader.Upgrade(w, r, nil)
		c.EnableWriteCompression(true)
		go handleNewClient(c)
	}
}

func main() {
	var upgrader = websocket.Upgrader{EnableCompression: true}
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", WsStart(upgrader))
	//http.ListenAndServe(":8000", mux)

	//router := httprouter.New()
	//router.GET("/bchatws", WsStart(upgrader))
	log.Println("Server started on:", HUB_ADDR)
	log.Println(http.ListenAndServeTLS(HUB_ADDR, CERT_PEM, KEY_PEM, mux))

}
