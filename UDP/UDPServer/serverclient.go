package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

func main() {
	var saddr net.UDPAddr
	saddr.IP = net.ParseIP("127.0.0.1")
	saddr.Port = 8888
	connection, _ := net.ListenUDP("udp", &saddr)
	for i := 0; i < runtime.NumCPU(); i++ {
		fmt.Println(runtime.NumCPU())
		go handleConnect(connection)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		if strings.Compare(msg, "close server\n") == 0 {
			break
		}
	}
}
func handleConnect(connection *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		memset(buffer)
		_, client, _ := connection.ReadFromUDP(buffer)
		connection.WriteToUDP(buffer, client)
		fmt.Println(client.Port, client.IP, string(buffer))
	}
}
func memset(des []byte) {
	for i := range des {
		des[i] = byte(0)
	}
}
