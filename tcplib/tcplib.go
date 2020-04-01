package tcplib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/yangyouwei/ws-tcp-monitor/loglib"
	"log"
	"net"
	"strings"
)

var key = "dAr1HMcHYltalEJoDdeUexMM7JWK5FQAilIAtGUkuVkfbI/aS4IH/JJ1vNlOybCW8KKvsgFadO6dxiGi97g83tGpPWqkcP56xvzJtEgC8XQGOq0OZS8MiavbMNrawtNHmgDEUMUaRAwpb7sPag161hEMtMPQ/Bkwdde8plzR4a8="
var body = "5X79PPxj5rjsKqk="
var macid = "2b344119bab20a66753136d7af65eca7"
const secwd  = "eyJzdGF0dXMiOjIwMH0="
var logtofile *log.Logger


func MakePackage() []byte {
	str := macid + key + body
	//uint16 big endian
	TotalLen := turtobyte(10 + len(str))
	KeyLen := turtobyte(len(key))
	MacidLen := turtobyte(len(macid))
	ReqId := turtobyte(1234)
	Cmd := turtobyte(1)

	//make head []byte
	headslice := [][]byte{TotalLen,KeyLen,MacidLen,ReqId,Cmd}
	Head := bytes.Join(headslice,[]byte{})
	//make tail []byte
	tail := []byte(str)
	//make full package
	fullslice := [][]byte{Head,tail}
	fullpack := bytes.Join(fullslice,[]byte{})
	//send package to server
	return fullpack
}

func TcpConnect(s string,m []byte) bool {
        var tcpAddr *net.TCPAddr
	logtofile = loglib.Logtofile
	tcpAddr,_ = net.ResolveTCPAddr("tcp",s)
	conn,err := net.DialTCP("tcp",nil,tcpAddr)

	if err!=nil {
		fmt.Println("Client connect error ! " + err.Error())
		logtofile.Println(tcpAddr,": tcp connect is fail.")
		return false
	}
	defer conn.Close()
	_, err = conn.Write(m)

	if err != nil {
		logtofile.Println(err)
	}

	var buf = make([]byte, 32)
	n, err := conn.Read(buf)
	if err != nil {
		logtofile.Println("read error:", err)
		logtofile.Println("tcp connect is fail.")
		return false
	} else {
		//logtofile.Printf("read %v bytes, content is %s\n", n, string(buf[:n]))
		//fmt.Println(string(buf[:n]))
		if strings.HasSuffix(string(buf[:n]),"5X7jO/p+9+n0IuYanKs/Zv") {
			fmt.Println(tcpAddr,": Tcp Connect Success!")
			//logtofile.Println("Tcp Connect Success!")
			return true
		}
	}
	logtofile.Println("tcp connect is fail.")
	return false
}

func turtobyte(lenint int) []byte {
	lenuint16 := uint16(lenint)
	bytebuffer := bytes.NewBuffer(make([]byte,0,1024))
	binary.Write(bytebuffer,binary.BigEndian,lenuint16)
	return  bytebuffer.Next(2)
}
