package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

func main() {
	var saddr net.TCPAddr
	saddr.IP = net.ParseIP("127.0.0.1")
	saddr.Port = 8000
	listen, _ := net.ListenTCP("tcp", &saddr)
	for {
		conect, _ := listen.AcceptTCP()
		go send(conect)
		go recv(conect)
	}
}
func send(connect net.Conn) {
	for {
		time.Sleep(1 * time.Second)
		connect.Write([]byte("Toang\n"))
	}
}
func recv(connect net.Conn) {
	for {
		reader := bufio.NewReader(connect)
		msg, _ := reader.ReadString('\n')
		log.Println(msg)
	}
}
