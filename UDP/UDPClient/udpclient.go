package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	var saddr net.UDPAddr
	saddr.Port = 8888
	saddr.IP = net.ParseIP("127.0.0.1")
	server, _ := net.DialUDP("udp", nil, &saddr)
	fmt.Println("running")
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		server.Write([]byte(msg))
	}
}
