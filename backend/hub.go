package main

import (
	"fmt"
	"github.com/engineerbeard/barrenschat/httpscerts"
	b "github.com/engineerbeard/barrenschat/shared"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const CERT_PEM = "cert.pem"
const KEY_PEM = "key.pem"


var serverObj ServerStruct

func init() {
	rand.Seed(time.Now().UnixNano())
	serverObj.Clients = make(map[string][]b.BChatClient)
	serverObj.mu = sync.Mutex{}
	serverObj.AddRoom(b.MAIN_ROOM)
}

func main() {
	HUB_ADDR := ":" + os.Getenv("PORT")
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

	log.Println("Server started on:", HUB_ADDR)
	log.Println(http.ListenAndServeTLS(HUB_ADDR, CERT_PEM, KEY_PEM, mux))
	//log.Println(http.ListenAndServe(HUB_ADDR, mux))
}

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

	// Move client to new room if the room exists
	if _, roomExists := s.Clients[bMessage.Room]; roomExists {
		s.Clients[bMessage.Room] = append(s.Clients[bMessage.Room], s.Clients[room][idx])
	} else { // Make a new room and move the client to it
		s.Clients[bMessage.Room] = []b.BChatClient{s.Clients[room][idx]}
	}

	//Remove client from original room
	s.Clients[room] = append(s.Clients[room][:idx], s.Clients[room][idx+1:]...)
	s.mu.Unlock()

	// Update clients
	// Broadcast to new room
	bMessage.RoomData = s.GetRoomsString()
	bMessage.OnlineData = s.GetNamesInRoom(bMessage.Room)
	bMessage.Payload = fmt.Sprint(bMessage.Name, " joined the room.")

	// Broadcast to old room
	s.BroadcastMessage(bMessage.Room, bMessage)
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
		keys = fmt.Sprint(keys, k, "\n")
		//keys = keys + k + "\n"
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
	log.Println(m)
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

func handleNewClient(c *websocket.Conn) {
	defer c.Close()
	var clientId string

	BClient := b.BChatClient{WsConn: c}
	bMessage, _ := BClient.ReadMessage()
	BClient.Name = bMessage.Payload
	BClient.Room = b.MAIN_ROOM
	BClient.Uid = bMessage.Uid

	serverObj.AddClient(&BClient)
	BClient.SendMessage(b.BMessage{
		MsgType:    b.B_CONNECT,
		Name:       bMessage.Payload,
		Room:       b.MAIN_ROOM,
		Payload:    "Welcome to BarrensChat\nCommands:\n /name <name>\n /room <room>\n",
		TimeStamp:  time.Now(),
		OnlineData: serverObj.GetNamesInRoom(b.MAIN_ROOM),
		RoomData:   serverObj.GetRoomsString(),
	})

	clientId = bMessage.Uid
	serverObj.BroadcastMessage(b.MAIN_ROOM, b.BMessage{
		MsgType:    b.B_CONNECT,
		Name:       bMessage.Payload,
		Room:       b.MAIN_ROOM,
		Payload:    "New Connection!",
		TimeStamp:  time.Now(),
		OnlineData: serverObj.GetNamesInRoom(b.MAIN_ROOM),
		RoomData:   serverObj.GetRoomsString(),
	})

	for {
		bMessage, err := BClient.ReadMessage()
		if err == nil { // Process message

			room, idx, _ := serverObj.FindClient(clientId)
			bMessage.OnlineData = serverObj.GetNamesInRoom(room)
			if bMessage.MsgType == b.B_MESSAGE {
				serverObj.BroadcastMessage(room, bMessage)
			} else if bMessage.MsgType == b.B_NAMECHANGE {
				serverObj.ChangeClientName(room, idx, bMessage.Name)
				bMessage.OnlineData = serverObj.GetNamesInRoom(room)
				serverObj.BroadcastMessage(room, bMessage)
			} else if bMessage.MsgType == b.B_ROOMCHANGE {
				serverObj.ChangeClientRoom(room, idx, bMessage)
			}
		} else { // Clean up
			BClient.Close()
			room, idx, name := serverObj.FindClient(clientId)
			serverObj.RemoveClient(room, idx)
			serverObj.BroadcastMessage(room, b.BMessage{
				MsgType:    b.B_DISCONNECT,
				TimeStamp:  time.Now(),
				Payload:    name + " Disconnected",
				OnlineData: serverObj.GetNamesInRoom(room),
				RoomData:   serverObj.GetRoomsString(),
			})
			log.Println(err) // Connection is over
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
