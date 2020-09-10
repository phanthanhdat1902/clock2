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

func main() {
	var saddr net.UDPAddr
	saddr.Port = 8888
	saddr.IP = net.ParseIP("192.168.1.7")
	server, _ := net.DialUDP("udp", nil, &saddr)
	fmt.Println("client running")
	for i := 0; i < 32; i++ {
		go recv(server)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		msg = editMsg(msg)
		if msg != "0" {
			msg += "\n"
			start := time.Now()
			//server.Write([]byte(msg))
			//reader = bufio.NewReader(server)
			//msg, _ = reader.ReadString('\n')
			//fmt.Println(msg)
			//fmt.Println("number C",numberC)
			for i := 0; i < runtime.NumCPU(); i++ {
				wg.Add(1)
				go func() {
					//fmt.Println(msg)
					start1 := time.Now().Second()
					for {
						server.Write([]byte(msg))
						numberC++
						fmt.Println("number C", numberC)
						end1 := time.Now().Second()
						if (end1 - start1) > 1 {
							break
						}
					}
					defer wg.Done()
				}()
			}
			wg.Wait()
			end := time.Now()
			PrintMemUsage()

			// Force GC to clear up, should see a memory drop
			runtime.GC()
			PrintMemUsage()
			fmt.Println(end, start, "\n-------------\n")
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
		fmt.Println(numberR)
		fmt.Println(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
