package main

import (
	"encoding/hex"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/libp2p/go-reuseport"
	"net"
	"runtime"
	"strconv"
	"strings"
)

var number = 0
var numberR = 0
var arrayCus []interface{}

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
	mongoSession := OpenDatabase()
	mongoSession.SetMode(mgo.Monotonic, true)
	runtime.GOMAXPROCS(100)
	for i := 0; i < 128; i++ {
		go decode()
		go workingDb(mongoSession)
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
		IP:   net.ParseIP("127.0.0.1"),
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
			dataWork.cus = cus
			dataWork.CMD = CMD
			databaseQueue <- dataWork
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
			dataWork.cus = cus
			dataWork.CMD = CMD
			databaseQueue <- dataWork
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
func OpenDatabase() *mgo.Session {
	session, _ := mgo.Dial("127.0.0.1:27017")
	return session
}

/*working db*/
func workingDb(mongoSession *mgo.Session) {
	sessionCopy := mongoSession.Copy()
	collection := sessionCopy.DB("my_exam").C("customer")
	for work := range databaseQueue {
		var resWork response
		resWork.err = "0"
		resWork.client = work.client
		resWork.CMD = work.CMD
		responseQueue <- resWork
		john := work.cus
		arrayCus = append(arrayCus, john)
		if len(arrayCus)%1000 == 0 {
			e := collection.Insert(arrayCus)
			if e != nil {
				fmt.Println(e)
			}
		}
		//}
	}
}

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
