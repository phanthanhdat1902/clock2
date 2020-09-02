package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
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
		result := handleMsg(string(buffer))
		connection.WriteToUDP([]byte(result), client)
		fmt.Println(client.Port, client.IP, string(buffer))
	}
}
func memset(des []byte) {
	for i := range des {
		des[i] = byte(0)
	}
}
func handleMsg(msg string) string {
	var result string
	//get CMD_MSISDN
	CMD_MSISDN := msg[:6]
	fmt.Println(CMD_MSISDN)
	//get IMSI
	lenIMSI, _ := strconv.Atoi(string(msg[6]))
	if lenIMSI > 11 {
		lenIMSI, _ = strconv.Atoi(hex.EncodeToString([]byte(strconv.Itoa(lenIMSI))))
	}
	index := 7 + (lenIMSI+1)/2
	IMSI := msg[7:index]
	fmt.Println(lenIMSI, IMSI)
	//getName
	lenName, _ := strconv.Atoi(string(msg[index]))
	if lenName == 0 {
		lenName, _ = strconv.Atoi(hex.EncodeToString([]byte(string(msg[index]))))
	}
	name := msg[index+1 : index+lenName+1]
	index += lenName + 1
	fmt.Println(lenName, name)
	//get CMT
	lenCMT, _ := strconv.Atoi(string(msg[index]))
	if lenCMT == 0 {
		lenCMT, _ = strconv.Atoi(hex.EncodeToString([]byte(string(msg[index]))))
	}
	CMT := msg[index+1 : index+lenCMT+1]
	fmt.Println(lenCMT, CMT)
	index += lenCMT + 1
	//get Birthday
	lenBirthday, _ := strconv.Atoi(string(msg[index]))
	if lenBirthday == 0 {
		lenBirthday, _ = strconv.Atoi(hex.EncodeToString([]byte(string(msg[index]))))
	}
	birthday := msg[index+1 : index+lenBirthday+1]
	fmt.Println(lenBirthday, birthday)
	return result
}
