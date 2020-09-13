package main

import (
	"encoding/hex"
	"fmt"
	"github.com/libp2p/go-reuseport"
	"net"
	"runtime"
	"strconv"
	"strings"
)
import _ "github.com/libp2p/go-reuseport"

import _ "github.com/go-sql-driver/mysql"

var number = 0
var numberR = 0

//const connect database
const (
	DB_HOST = "tcp(127.0.0.1:3306)"
	DB_NAME = /*name database*/ "my_exam"
	DB_USER = /*"user"*/ "ptd"
	DB_PASS = /*"pass"*/ "anh123asd"
)

/*model*/
type customer struct {
	MSISDN   string
	IMSI     string
	Name     string
	CMT      string
	Birthday string
}

type works struct {
	client net.Addr
	msg    string
}

type response struct {
	client net.Addr
	err    string
	CMD    int
}

type databaseWorks struct {
	client net.Addr
	cus    customer
	CMD    int
}

var numberC = 0
var numberD = 0
var done chan bool
var decodeQueue = make(chan works, 500)
var databaseQueue = make(chan databaseWorks, 500)
var responseQueue = make(chan response, 500)

func main() {
	runtime.GOMAXPROCS(100)
	for i := 0; i < 128; i++ {
		go decode()
		//go workingDb()
	}
	for i := 0; i < 64; i++ {
		go listening()
	}
	fmt.Println("server running")
	<-done
}

func listening() {
	addr := net.UDPAddr{
		Port: 8888,
		IP:   net.ParseIP("192.168.1.150"),
	}

	connection, err := reuseport.ListenPacket("udp", addr.String())
	for i := 0; i < 12; i++ {
		go responseClient(connection)
	}
	if err != nil {
		panic(err)
	}
	for {
		var work works
		buffer := make([]byte, 200)
		_, client, _ := connection.ReadFrom(buffer)
		work.client = client
		work.msg = string(buffer)
		decodeQueue <- work
		numberC++
		fmt.Println("number : ", numberC)
	}
}

func decode() {
	var resWork response
	for work := range decodeQueue {
		number++
		fmt.Println(number)
		var dataWork databaseWorks
		msg := work.msg
		dataWork.client = work.client
		var cus customer
		//get CMD_MSISDN
		CMD_MSISDN := msg[:6]
		CMD_MSISDN = hex.EncodeToString([]byte(CMD_MSISDN))
		CMD, _ := strconv.Atoi(string(CMD_MSISDN[0]))
		MSISDN := CMD_MSISDN[1:]
		cus.MSISDN = MSISDN
		//get IMSI
		lenIMSI, _ := strconv.Atoi(hex.EncodeToString([]byte(string(msg[6]))))
		lenIMSI -= 10
		index := 7 + (lenIMSI+1)/2
		IMSI := msg[7:index]
		IMSI = hex.EncodeToString([]byte(IMSI))
		IMSI = IMSI[:lenIMSI]
		cus.IMSI = IMSI
		if CMD == 3 {
			//dataWork.cus = cus
			//dataWork.CMD = CMD
			//databaseQueue <- dataWork
			return
		} else {
			//getName
			lenName, _ := strconv.Atoi(hex.EncodeToString([]byte(string(msg[index]))))
			lenName -= 10
			name := msg[index+1 : index+lenName+1]
			index += lenName + 1
			cus.Name = name
			//get CMT
			lenCMT, _ := strconv.Atoi(hex.EncodeToString([]byte(string(msg[index]))))
			lenCMT -= 10
			CMT := msg[index+1 : index+lenCMT+1]
			cus.CMT = CMT
			index += lenCMT + 1
			//get Birthday
			lenBirthday, _ := strconv.Atoi(hex.EncodeToString([]byte(string(msg[index]))))
			lenBirthday -= 10
			birthday := msg[index+1 : index+lenBirthday+1]
			cus.Birthday = birthday
			resWork.client = work.client
			resWork.err = "0"
			resWork.CMD = 1
			responseQueue <- resWork
			//dataWork.cus = cus
			//dataWork.CMD = CMD
			//databaseQueue <- dataWork
		}
	}
}

/*open*/
//func OpenDB() (*sql.DB, error) {
//	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME))
//	if err != nil {
//		fmt.Println(err)
//		return db, err
//	}
//	return db, nil
//}

/*open db mongo*/
//func OpenDatabase()  *mongo.Collection{
//	// Set client options
//	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
//
//	// Connect to MongoDB
//	client, e := mongo.Connect(context.TODO(), clientOptions)
//	if e!=nil{
//		fmt.Println(e)
//	}
//	// Check the connection
//	e = client.Ping(context.TODO(), nil)
//	if e!=nil{
//		fmt.Println(e)
//	}
//	collection := client.Database("my_exam").Collection("customer")
//	return collection
//}
/*working db*/
//func workingDb()  {
//	collection:=OpenDatabase()
//	for work := range databaseQueue {
//		var resWork response
//		john := work.cus
//		resWork.client = work.client
//		resWork.CMD = work.CMD
//		_, e := collection.InsertOne(context.Background(), john)
//		if e != nil {
//			fmt.Println(e)
//		} else {
//			numberD++
//			fmt.Println("numberD: ", numberD)
//		}
//		resWork.err = "0"
//		responseQueue <- resWork
//		continue
//	}
//}

func responseClient(connection net.PacketConn) {
	for work := range responseQueue {
		if strings.Compare(work.err, "0") == 0 || strings.Compare(work.err, "") == 0 {
			connection.WriteTo([]byte("200 OK\n"), work.client)
		} else {
			fmt.Println("a")
			connection.WriteTo([]byte("400 error: "+work.err+"\n"), work.client)
		}
		numberR++
		fmt.Println(numberR)
	}
}
