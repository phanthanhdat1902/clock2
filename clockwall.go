package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) > 1 {
		for i:=1;i<len(os.Args);i++{
			//Lay ra tung server, tach chuoi cau truc: nameserver=port
			server:=strings.Split(os.Args[i],"=")
			//truyen vao ten va port, chay moi server tren 1 goroutines
			go runServer(server[0],server[1])
		}
		for  {

		}
	}else{
		fmt.Println("Hay nhap vao so cong tuong ung")
	}
}
func runServer(nameServer string,port string)  {
	listener, err := net.Listen("tcp", port)
	if err!=nil{
		log.Fatal(err)
	}else{
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err) // e.g., connection aborted
				continue
			}
			go handleOneConn(nameServer,conn) // Xu ly nhieu ket noi tai 1 thoi diem
		}
	}
}
func handleOneConn(nameServer string,c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, nameServer+": "+time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}