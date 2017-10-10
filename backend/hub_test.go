package main

import (
	"math/rand"
	"net/http"
	"testing"

	"strings"

	"net/http/httptest"

	"net/url"

	bs "github.com/engineerbeard/barrenschat/shared"
	//"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var rndUid, rndName, rndRoom string
var names string

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func init() {
	log.SetOutput(ioutil.Discard)
	serverObj.Clients = make(map[string][]bs.BChatClient)
	rndUid = RandStringRunes(32)
	rndName = RandStringRunes(32)
	rndRoom = RandStringRunes(32)
	serverObj.Clients["TEST"] = append(serverObj.Clients["TEST"], bs.BChatClient{Uid: rndUid, Name: rndName, Room: rndRoom})
	for j := 0; j < 15; j++ {
		q := RandStringRunes(10)
		serverObj.Clients[q] = []bs.BChatClient{}
		for j := 0; j < 10; j++ {
			n := RandStringRunes(10)
			serverObj.Clients[q] = append(serverObj.Clients[q], bs.BChatClient{Uid: RandStringRunes(32), Name: n, Room: RandStringRunes(10)})
		}
	}
}
func TestFindClient(t *testing.T) {
	room, _, name := serverObj.FindClient(rndUid)
	if room != rndRoom {
		t.Fatal("Incorrect Room")
	}
	if name != rndName {
		t.Fatal("Incorrect Uid")
	}
}
func BenchmarkFindClient(b *testing.B) {
	for n := 0; n < b.N; n++ {
		room, _, name := serverObj.FindClient(rndUid)
		if room != rndRoom {
			b.Fatal("Incorrect Room")
		}
		if name != rndName {
			b.Fatal("Incorrect Uid")
		}
	}
}
func TestFindClientFail(t *testing.T) {
	room, idx, name := serverObj.FindClient(RandStringRunes(10))
	if room != "" {
		t.Fatal("Incorrect Room")
	}
	if name != "" {
		t.Fatal("Incorrect Uid")
	}
	if idx != 0 {
		t.Fatal("Incorrect Idx")
	}
}
func BenchmarkFindClientFail(b *testing.B) {
	for n := 0; n < b.N; n++ {
		room, idx, name := serverObj.FindClient("asdf")
		if room != "" {
			b.Fatal("Incorrect Room")
		}
		if name != "" {
			b.Fatal("Incorrect Uid")
		}
		if idx != 0 {
			b.Fatal("Incorrect Idx")
		}
	}
}
func TestGetNamesInRoom(t *testing.T) {
	g := serverObj.GetNamesInRoom("TEST")
	g = strings.Replace(g, "\n", "", -1)
	if g != rndName {
		t.Fatalf("Names do not match")
	}
}
func TestGetNamesInRoomFail(t *testing.T) {
	g := serverObj.GetNamesInRoom("DNE")
	g = strings.Replace(g, "\n", "", -1)
	if g == rndName {
		t.Fatalf("Names do not match")
	}
}

func TestWsStart(t *testing.T) {
	uu := RandStringRunes(32)
	r := RandStringRunes(32)
	c, err := getClient(bs.MAIN_ROOM)
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.SendMessage(bs.BMessage{MsgType: bs.B_MESSAGE, Uid: uu, Payload: r})
	bMessage, err := c.ReadMessage()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if bMessage.Payload != r {
		t.Fatalf("Msg are not equal")
	}

	c.Close()

}
func TestBroadcastMessageRoomChange(t *testing.T) {
	uu := RandStringRunes(32)
	r := RandStringRunes(32)
	n := RandStringRunes(32)
	rr := RandStringRunes(32)
	c, err := getClient(bs.MAIN_ROOM)
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.SendMessage(bs.BMessage{MsgType: bs.B_ROOMCHANGE, Uid: uu, Payload: r, Name: n, Room: rr})
	bMessage, err := c.ReadMessage()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if bMessage.Name != n {
		t.Fatalf("Name are not equal")
	}
	if bMessage.Room != rr {
		t.Fatalf("Name are not equal")
	}
	if !strings.Contains(bMessage.RoomData, rr) {
		t.Fatalf("Doesnt contain room name")
	}

	c.Close()
}
func TestBroadcastMessageNameChange(t *testing.T) {
	uu := RandStringRunes(32)
	r := RandStringRunes(32)
	n := RandStringRunes(32)
	c, err := getClient(bs.MAIN_ROOM)
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.SendMessage(bs.BMessage{MsgType: bs.B_NAMECHANGE, Uid: uu, Payload: r, Name: n})
	bMessage, err := c.ReadMessage()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if bMessage.Payload != r {
		t.Fatalf("Msg are not equal")
	}
	if bMessage.Name != n {
		t.Fatalf("Name are not equal")
	}
	c.Close()
	//serverObj = make(map[string][]BChatClient)
}

//func TestConnect10Clients(t *testing.T) {
//	clientList := []bs.BChatClient{}
//	for n := 0; n < 10; n++ {
//		f, err := getClient(bs.MAIN_ROOM)
//		if err != nil {
//			t.Fatalf(err.Error())
//		}
//		clientList = append(clientList, f)
//		go func() {
//			for {
//				f.ReadMessage()
//			}
//		}()
//		//uu := RandStringRunes(32)
//		//r := RandStringRunes(32)
//		////c, err := getClient()
//		//if err != nil {
//		//	t.Fatalf(err.Error())
//		//}
//		//f.SendMessage(bs.BMessage{MsgType:bs.B_MESSAGE, Uid:uu, Payload:r})
//		//c.ReadMessage()
//	}
//	for _, clients := range clientList {
//		clients.Close()
//	}
//}

//func message(i int, b *testing.B) {
//	var bMessage bs.BMessage
//	conn, _ := getWsConn()
//	uu := RandStringRunes(32)
//	BClient := BChatClient{WsConn: conn, Uid: RandStringRunes(32), Name: "Anon", Room: bs.MAIN_ROOM}
//	BClient.SendMessage(bs.BMessage{MsgType: bs.B_CONNECT, Uid: uu, Payload: BClient.Name})
//	bMessage, _ = BClient.ReadMessage()
//	if !strings.Contains(bMessage.Payload, "Welcome") {
//		b.Fatalf("Does not contain welcome", bMessage.Payload)
//	}
//	bMessage, _ = BClient.ReadMessage()
//	if !strings.Contains(bMessage.Payload, "Connect") {
//		b.Fatalf("Does not contain welcome", bMessage)
//	}
//	for n := 0; n < b.N; n++ {
//		r := RandStringRunes(i)
//		BClient.SendMessage(bs.BMessage{MsgType:bs.B_MESSAGE, Uid:uu, Payload:r})
//		msg, _ := BClient.ReadMessage()
//		if msg.Payload != r {
//			b.Fatal("Bad msg")
//		}
//	}
//	BClient.Close()
//}
//func BenchmarkMessages(b *testing.B) {
//	message(10, b)
//}

func getClient(room string) (bs.BChatClient, error) {
	var bMessage bs.BMessage
	var upgrader = websocket.Upgrader{EnableCompression: true}
	var BClient bs.BChatClient
	srv := httptest.NewServer((WsStart(upgrader)))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return BClient, err
	}
	conn.EnableWriteCompression(true)
	BClient = bs.BChatClient{WsConn: conn, Uid: RandStringRunes(32), Name: RandStringRunes(10), Room: room}
	if err != nil {
		return BClient, err
	}

	BClient.SendMessage(bs.BMessage{MsgType: bs.B_CONNECT, Uid: RandStringRunes(32), Payload: BClient.Name})
	bMessage, err = BClient.ReadMessage()

	if !strings.Contains(bMessage.Payload, "Welcome") {
		return BClient, err
	}
	bMessage, err = BClient.ReadMessage()
	if !strings.Contains(bMessage.Payload, "Connection!") {
		return BClient, err
	}
	return BClient, nil
}

func getRequestWithGET(t testing.TB, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
