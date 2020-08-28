package main

import (
	"bufio"
	"fmt"
	"net"
)

var(
	conns []net.Conn
	mess=make(chan string)
	connect=make(chan net.Conn)
	connectLeaves=make(chan net.Conn)
)
func main() {
	listen,err:=net.Listen("tcp","localhost:8888")
	if err!=nil{
		fmt.Println(err)
	}
	go func() {
		for{
			conn,err:=listen.Accept()
			if err!=nil{
				fmt.Println(err)
				break
			}
			conns=append(conns,conn)
			connect<-conn
		}
	}()
	for{
		select {
			case conn:=<-connect:
				go recv(conn)
			//case conn:=<-connectLeaves:
			//	fmt.Println(conn)
			//	fmt.Println("Thoat")
			case msg:=<-mess:
				fmt.Println(msg)
		}
	}
}
func recv(connect net.Conn)  {
	connect.Write([]byte("Hay nhap ten cua ban :D :"))
	reader:=bufio.NewReader(connect)
	name,_:=reader.ReadString('\n')
	name=name[:len(name)-1]
	for{
		reader:=bufio.NewReader(connect)
		msg,err:=reader.ReadString('\n')
		if err!=nil{
			break
		}
		mess<-msg
		fmt.Println(msg)
		sendAll(connect,name+":"+msg)
	}
	sendAll(connect,name+" da out")
	removeConn(connect)
	connect.Close()
}
func removeConn(connect net.Conn)  {

}
func sendAll(connect net.Conn,msg string)  {
	for i:= range conns{
		if conns[i]!=connect{
			conns[i].Write([]byte (msg))
		}
	}
}