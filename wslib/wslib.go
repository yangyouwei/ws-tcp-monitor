package wslib

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"github.com/yangyouwei/ws-tcp-monitor/loglib"
	"golang.org/x/net/websocket"
	"log"
)

const secwd  = "eyJzdGF0dXMiOjIwMH0="
var logtofile *log.Logger

type WSMessages struct {
	Msgid int `json:"msgid"`
}
func WebSocket(wsmessages WSMessages,wsurl string) bool {
	logtofile = loglib.Logtofile
	var origin = "http://localhost.com/"
	ws, err := websocket.Dial(wsurl, "", origin)
	if err != nil {
		fmt.Println(err)
		logtofile.Println(wsurl,": websocket connect fail.")
		return false
	}
	defer ws.Close()
	wsmj,err := json.Marshal(wsmessages)
	if err !=nil {
		logtofile.Println(err)
		fmt.Println("Websocket Connect Fail!")
		return false
	}
	        enwsmeg := base64.StdEncoding.EncodeToString(wsmj)
        //fmt.Println(enwsmeg)
        _, err = ws.Write([]byte(enwsmeg))
	if err != nil {
		logtofile.Println(err)
		fmt.Println("Websocket Connect Fail!")
		return false
	}

	data := make([]byte, 20)
	_, err = ws.Read(data)
	if err != nil {
		logtofile.Println(err)
		fmt.Println("Websocket Connect Fail!")
		logtofile.Println("websocket read messages fail.")
		return false
	}

	if string(data) == secwd {
		//fmt.Println(string(data))
		//logtofile.Println("Websocket Connect Secuess!")
		fmt.Println(wsurl,": Websocket Connect Secuess!")
		return true
	}else {
		fmt.Println("Websocket Connect Fail!")
		return false
	}
}
