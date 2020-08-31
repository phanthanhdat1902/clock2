package main

import (
	"fmt"
	"net"
	"runtime"
)

func main() {
	var saddr net.UDPAddr
	saddr.IP = net.ParseIP("127.0.0.1")
	saddr.Port = 8888
	connection, _ := net.ListenUDP("udp", &saddr)
	for i := 0; i < runtime.NumCPU(); i++ {
		go handleConnect(connection)
	}
	for {

	}
}
func handleConnect(connection *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		_, client, _ := connection.ReadFromUDP(buffer)
		fmt.Println(client.Port, client.IP, string(buffer))
	}
}
