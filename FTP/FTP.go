package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
)

func main() {
	listen,err:=net.Listen("tcp","localhost:8000")
	if err!=nil{
		fmt.Println(err)
	}
	for{
		connect,err:=listen.Accept()
		if err!=nil{
			fmt.Println(err)
			break
		}
		go recv(connect)
	}
	//out, err := exec.Command("cmd","/C", "dir").Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf(" %s\n", out)
}
func recv(connect net.Conn)  {
	//lay thu muc hien tai
	dir, err := os.Getwd()
	dir="[server]"+dir+">"
	if err != nil {
		fmt.Println(err)
	}
	connect.Write([]byte ("Chao mung ban den voi FTP server\n"+dir))
	for{
		reader:=bufio.NewReader(connect)
		msg,err:=reader.ReadString('\n')
		if err!=nil{
			fmt.Println(err)
			break
		}
		out, err := exec.Command("cmd","/C", msg).Output()
		connect.Write(out)
		connect.Write([]byte (dir))
	}
}
