package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
)

func main() {
	_, err := hex.DecodeString("8")
	fmt.Println(err)
	//var saddr net.TCPAddr
	//saddr.IP = net.ParseIP("127.0.0.1")
	//saddr.Port = 8000
	//connect, _ := net.DialTCP("tcp", nil, &saddr)
	//go recv(connect)
	//connect.CloseWrite()
	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	msg, _ := reader.ReadString('\n')
	//	connect.Write([]byte(msg))
	//}
}
func recv(connect net.Conn) {
	for {
		reader := bufio.NewReader(connect)
		msg, _ := reader.ReadString('\n')
		fmt.Println(msg)
	}
}
