package main


import (
	"math/rand"
	"testing"
	"net/http"

)
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var rndUid, rndName, rndRoom string
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func init() {
	connectedClients = make(map[string][]BChatClient)
	rndUid = RandStringRunes(32)
	rndName = RandStringRunes(32)
	rndRoom = RandStringRunes(32)
	connectedClients["TEST"] = append(connectedClients["TEST"] ,BChatClient{Uid:rndUid, Name:rndName, Room:rndRoom})
	for j := 0; j < 2; j++ {
		q := RandStringRunes(10)
		connectedClients[q] = []BChatClient{}
		for j := 0; j < 10; j++ {
			connectedClients[q] = append(connectedClients[q] ,BChatClient{Uid:RandStringRunes(32), Name:RandStringRunes(10), Room:RandStringRunes(10)})
		}
	}
}

func TestFindClient(t *testing.T) {
	room, _, name := FindClient(rndUid)
	if room != rndRoom {
		t.Fatal("Incorrect Room")
	}
	if name != rndName {
		t.Fatal("Incorrect Uid")
	}
}

func TestFindClientFail(t *testing.T) {
	room, idx, name := FindClient(RandStringRunes(10))
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

func TestGetNamesInRoom(t *testing.T) {

}

//func TestWsStart(t *testing.T) {
//	var newBox bchatcommon.BMessage
//
//	r := getRequestWithGET(t, "/uncompressedget")
//
//	rw := httptest.NewRecorder()
//	fn := WsStart("s")
//	fn(rw, r, httprouter.Params{})
//
//	z := bytes.NewReader(rw.Body.Bytes())
//
//	err := json.NewDecoder(z).Decode(&newBox)
//
//	if err != nil {
//
//		t.Fatal("Error decodeing json", err.Error())
//
//	}
//
//}

func getRequestWithGET(t testing.TB, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req

}
