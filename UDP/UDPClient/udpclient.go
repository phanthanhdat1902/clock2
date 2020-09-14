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
	"sync"
	"time"
)

var wg sync.WaitGroup
var numberC = 0
var numberR = 0
var done chan bool
var start int
var end int

func main() {
	var saddr net.UDPAddr
	saddr.Port = 8888
	saddr.IP = net.ParseIP("127.0.0.1")
	server, _ := net.DialUDP("udp", nil, &saddr)
	fmt.Println("client running")
	for i := 0; i < 32; i++ {
		go recv(server)
	}
	start = time.Now().Second()
	go input(server)
	for {
		end = time.Now().Second()
		if (end - start) >= 10 {
			os.Exit(3)
		}
	}
	<-done
}
func input(server *net.UDPConn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		msg = editMsg(msg)
		if msg != "0" {
			msg += "\n"
			for i := 0; i < runtime.NumCPU(); i++ {
				wg.Add(1)
				go func() {
					//fmt.Println(msg)
					for {
						server.Write([]byte(msg))
						numberC++
						fmt.Println("number C", numberC)
					}
					wg.Done()
				}()
			}
			wg.Wait()
		} else {
			fmt.Println("Du lieu truyen vao khong dung dinh dang, hay phan tach cac truong bang dau cach")
		}

	}
}
func editMsg(msg string) string {
	var result string
	msg = msg[:len(msg)-1]
	msg = strings.TrimSpace(msg)
	trimMsg := strings.Split(msg, " ")
	//check CMD
	if len(trimMsg) == 3 {
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
		lenIMSI += 10
		temp, _ = hex.DecodeString(strconv.Itoa(lenIMSI))
		result += string(temp)
		temp, _ = hex.DecodeString(IMSI)
		result += string(temp)
		return result
	} else if len(trimMsg) == 6 {
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
		lenIMSI += 10
		temp, _ = hex.DecodeString(strconv.Itoa(lenIMSI))
		result += string(temp)
		temp, _ = hex.DecodeString(IMSI)
		result += string(temp)
		//check name
		name := trimMsg[3]
		lenName := len(name)
		lenName += 10
		temp, _ = hex.DecodeString(strconv.Itoa(lenName))
		result += string(temp)
		result += name
		//check CMT
		CMT := trimMsg[4]
		lenCMT := len(CMT)
		lenCMT += 10
		temp, _ = hex.DecodeString(strconv.Itoa(lenCMT))
		result += string(temp)
		//fmt.Println(temp)
		result += CMT
		//check birthday
		birthday := trimMsg[5]
		lenBirthday := len(birthday)
		lenBirthday += 10
		temp, _ = hex.DecodeString(strconv.Itoa(lenBirthday))
		result += string(temp)
		result += birthday
		return result
	} else {
		return "0"
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v Kib", bToKb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v Kib", bToKb(m.TotalAlloc))
	fmt.Printf("\tSys = %v KiB", bToKb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func recv(server *net.UDPConn) {
	for {
		reader := bufio.NewReader(server)
		res, err := reader.ReadString('\n')
		numberR++
		fmt.Println("numberR: ", numberR)
		fmt.Println(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
