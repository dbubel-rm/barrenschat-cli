package main

// language: go
// sudo: false
// matrix:
//   include:
//   - go: 1.x
//     env: LATEST=true
// before_install:
// - go get github.com/mitchellh/gox
// install:
// -
// script:
// - go get -t -v ./...
// -
// - go test -v -race ./... -bench=.
// - if [ "${LATEST}" = "true" ]; then gox -os="windows linux darwin" -arch="amd64" -verbose
//   ./...; fi
// deploy:
//   provider: releases
//   api_key:
//     secure: fQxUOsW7veYOzyAcxJmVrUjMDsdn/eXzw9KoEUXnkZ4auJLvfhbBhqnJV0NKI9/7sM/KylLVHgoLPuHv5Ne6KQ5XQBuLWIKAm7NkaLfNNFi1JOjYoM3jMCEDoqawaNAhtU5EkCUQEKLhcATC9GMMhR2vPn1Ak4aLOXSUXQuVHajt4oCvUK84OduM9FmHvS/mAwL9g0U/HkvTzHYxQMkiUkREc8kUrkbABh0WtrMWwPguYIfQjG79XdeFnAFsWYakodfZK76ao59O/GHh0BYrLpcazAaK5JU2FpvTS5RiwOOUbR7zn8HYUk+hn1ol5ZFZyGz5n82VmbT9gK43fH1s5dZg+KiJCa5ITtW8YNcAfc2Gb1lciooV7xzMMFBj3NIhG6XISrVQ507Qh2Qdkg5hVGxU/7vSewtadGxOC7Z2ILKLQq5ckFicDsh+ZEIJObdOCFWy0VCox4Ei+fT+SFS8EpDn1u4bv8jfpq/9MAu/HDeHUeA1v0xHT5K5ZCqtsMWskROrJMfY34l50nzytb2ZxRZ4B8Oe1/x3/cB8zsq5jSkV7ND2ca4a3iPnwFica4/EShuhswSbodZk5TbI5MTuDGcRvonFwzPZ+Bcxrq+rYs2Q1idaPE6AmYOgkiFkz2JuecGUKaAc/7dhxBjx7O4M2ODzS7X13RdJYkTtUEgdBv0=
//   file:
//   - frontend_windows_amd64.exe
//   - frontend_linux_amd64
//   - frontend_darwin_amd64
//   on:
//     branch: master

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"

	"math/rand"
	"time"

	//"net/url"

	"strings"
)

const CHATWINDOW = "CHATWINDOW"
const ONLINEWINDOW = "ONLINEWINDOW"
const ROOMWINDOW = "ROOMWINDOWs"

// type server struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
	// BClient = b.BChatClient{}

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	a := make([]rune, n)
	for i := range a {
		a[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(a)
}
func main() {
	// Setup Ws connection
	// u := url.URL{Scheme: "wss", Host: "damp-springs-83733.herokuapp.com", Path: "/bchatws"}
	// d := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	// c, _, err := d.Dial(u.String(), nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	//c.WriteJSON()

	//bhelpersb..BMessage{MsgType:B_CONNECT, Uid:RandStringRunes(32)}
	// BClient = b.BChatClient{WsConn: c, Uid: RandStringRunes(32), Name: "Anon", Room: b.MAIN_ROOM}
	// BClient.SendMessage(b.BMessage{MsgType: b.B_CONNECT, Uid: RandStringRunes(32), Payload: BClient.Name})

	// Setup CUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true

	g.SetManagerFunc(setLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, onEnterEvt()); err != nil {
		log.Panicln(err)
	}

	// go handleConnection(c, g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// func handleConnection(c *websocket.Conn, g *gocui.Gui) {
// 	var bMessage b.BMessage
// 	c.ReadJSON(&bMessage)
// 	processMsg(bMessage, g)
// 	for {
// 		err := c.ReadJSON(&bMessage)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		processMsg(bMessage, g)
// 	}
// }

// func processMsg(msg b.BMessage, g *gocui.Gui) {
// 	g.Update(func(g *gocui.Gui) error {

// 		v, _ := g.View(CHATWINDOW)
// 		if msg.MsgType == b.B_CONNECT || msg.MsgType == b.B_DISCONNECT || msg.MsgType == b.B_ROOMCHANGE {
// 			o, _ := g.View(ONLINEWINDOW)
// 			o.Clear()
// 			fmt.Fprint(o, msg.OnlineData)

// 			o, _ = g.View(ROOMWINDOW)
// 			o.Clear()
// 			fmt.Fprintf(o, msg.RoomData)
// 		}
// 		if msg.MsgType == b.B_NAMECHANGE {
// 			o, _ := g.View(ONLINEWINDOW)
// 			o.Clear()
// 			fmt.Fprint(o, msg.OnlineData)
// 		}

// 		fmt.Fprintln(v, fmt.Sprintf("\u001b[33m%s\u001b[0m (%s) %s", msg.TimeStamp.Format("2006-01-02 15:04"), msg.Name, msg.Payload))
// 		return nil
// 	})
// }

func setActiveView(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func onEnterEvt() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		buf := strings.Replace(v.Buffer(), "\n", "", -1)
		if len(buf) < 1 {
			v.SetCursor(0, 0)
			return nil
		}
		// msgType := b.B_MESSAGE
		// if strings.Contains(buf, "/name") && len(strings.SplitAfter(buf, "/name")) > 1 {
		// 	newName := strings.SplitAfter(buf, "/name")[1]
		// 	newName = strings.TrimSpace(newName)
		// 	buf = fmt.Sprintf("%s changed name to %s", BClient.Name, newName)
		// 	BClient.Name = newName
		// 	msgType = b.B_NAMECHANGE

		// } else if strings.Contains(buf, "/room") && len(strings.SplitAfter(buf, "/room")) > 1 {
		// 	msgType = b.B_ROOMCHANGE
		// 	newRoom := strings.SplitAfter(buf, "/room")[1]
		// 	newRoom = strings.TrimSpace(newRoom)
		// 	BClient.Room = newRoom
		// 	//buf = fmt.Sprintf("%s left room", BClient.Name)
		// }
		// err := c.WriteJSON(b.BMessage{
		// 	MsgType:   msgType,
		// 	TimeStamp: time.Now(),
		// 	Name:      BClient.Name,
		// 	Room:      BClient.Room,
		// 	Uid:       BClient.Uid,
		// 	Payload:   buf,
		// })
		// if err != nil {
		// 	log.Println(err)
		// }

		v.Clear()
		v.SetCursor(0, 0)
		return nil
	}

}
func setLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(ONLINEWINDOW, 0, 0, 20, 14); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Online"
		v.Autoscroll = true
		fmt.Fprintln(v, "")
	}

	if v, err := g.SetView(ROOMWINDOW, 0, 15, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Rooms"
		fmt.Fprintln(v, "")
	}

	if v, err := g.SetView("input", 21, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
		v.Title = "Type To Chat"
		v.Editable = true
		v.Wrap = true
		//fmt.Fprintf(v, "H")
	}

	if v, err := g.SetView(CHATWINDOW, 21, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
		v.Title = "BChats"
		v.Editable = false
		v.Wrap = true
	}

	if _, err := setActiveView(g, "input"); err != nil {
		return err
	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
