package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
		msg = editMsg(msg)
		msg += "\n"
		server.Write([]byte(msg))
	}
}

func editMsg(msg string) string {
	var result string
	msg = msg[:len(msg)-1]
	trimMsg := strings.Split(msg, " ")
	//check CMD
	if strings.Compare(trimMsg[0], "3") == 0 {

	} else {
		//join CMD with MSISDN
		CMD_MSISDN := trimMsg[0] + trimMsg[1]
		CMD_MSISDNConver, _ := hex.DecodeString(CMD_MSISDN)
		result = string(CMD_MSISDNConver)
		//check len(IMSI)
		IMSI := trimMsg[2]
		lenIMSI := len(IMSI)
		if lenIMSI%2 != 0 {
			IMSI = IMSI + "0"
		}
		var temp []byte
		if lenIMSI < 11 {
			result += strconv.Itoa(lenIMSI)
		} else {
			temp, _ = hex.DecodeString(strconv.Itoa(lenIMSI))
			result += string(temp)
		}
		temp, _ = hex.DecodeString(IMSI)
		result += string(temp)
		//check name
		name := trimMsg[3]
		lenName := len(name)
		if lenName < 11 {
			result += strconv.Itoa(lenName)
		} else {
			temp, _ = hex.DecodeString(strconv.Itoa(lenName))
			result += string(temp)
		}
		result += name
		//check CMT
		CMT := trimMsg[4]
		lenCMT := len(CMT)
		if lenCMT < 11 {
			result += strconv.Itoa(lenCMT)
		} else {
			temp, _ = hex.DecodeString(strconv.Itoa(lenName))
			result += string(temp)
		}
		fmt.Println(temp)
		result += CMT
		//check birthday
		birthday := trimMsg[5]
		lenBirthday := len(birthday)
		if lenBirthday < 11 {
			result += strconv.Itoa(lenBirthday)
		} else {
			temp, _ = hex.DecodeString(strconv.Itoa(lenBirthday))
			result += string(temp)
		}
		result += birthday
		return result
	}
	return result
}
