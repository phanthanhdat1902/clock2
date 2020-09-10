package main

import (
	"database/sql"
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
	name     string
	CMT      string
	birthday string
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
var done chan bool
var decodeQueue = make(chan works, 500)
var databaseQueue = make(chan databaseWorks, 500)
var responseQueue = make(chan response, 500)

func main() {
	runtime.GOMAXPROCS(100)
	//var saddr net.UDPAddr
	//saddr.IP = net.ParseIP("127.0.0.1")
	//saddr.Port = 8888
	//connection, _ := net.ListenUDP("udp", &saddr)
	//for i:=0;i<256;i++{
	// go func() {
	//    for j:=0;j<10;j++{
	//       go recv(connection)
	//    }
	// }()
	//}
	//for i:=0;i<64;i++{
	// go responseClient(connection)
	//}
	for i := 0; i < 128; i++ {
		go decode()
		go sendDatabase()
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
		IP:   net.ParseIP("192.168.1.7"),
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
			cus.name = name
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
			cus.birthday = birthday
			dataWork.cus = cus
			dataWork.CMD = CMD
			databaseQueue <- dataWork
		}
	}
}

/*open*/
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME))
	if err != nil {
		fmt.Println(err)
		return db, err
	}
	return db, nil
}

func sendDatabase() {
	db, err := OpenDB()
	for work := range databaseQueue {
		var resWork response
		getCus := work.cus
		resWork.client = work.client
		resWork.CMD = work.CMD
		if db != nil {
			if work.CMD == 1 {
				stmt, err := db.Prepare("insert into infor_customer values (?,?,?,?,?)")
				if err != nil {
					resWork.err = err.Error()
					responseQueue <- resWork
					fmt.Println("workingDtabase db.Prepare", err)
					continue
				} else {
					_, err := stmt.Exec(getCus.MSISDN, getCus.IMSI, getCus.name, getCus.CMT, getCus.birthday)
					if err != nil {
						resWork.err = err.Error()
						responseQueue <- resWork
						continue
					}
				}
			} else if work.CMD == 2 {
				var cus customer
				row := db.QueryRow("select * from infor_customer where MSISDN=" + getCus.MSISDN)
				err := row.Scan(&cus.MSISDN, &cus.IMSI, &cus.name, &cus.CMT, &cus.birthday)
				if err != nil {
					resWork.err = err.Error()
					responseQueue <- resWork
					continue
				}
				if getCus.birthday == "0" {
					getCus.birthday = cus.birthday
				}
				if getCus.CMT == "0" {
					getCus.CMT = cus.CMT
				}
				if getCus.name == "0" {
					getCus.name = cus.name
				}
				stmt, err := db.Prepare("update infor_customer set full_name=?,CMND=?,birthday=? where MSISDN=?")
				if err != nil {
					resWork.err = err.Error()
					responseQueue <- resWork
					continue
				} else {
					res, err := stmt.Exec(getCus.name, getCus.CMT, getCus.birthday, getCus.MSISDN)
					a, _ := res.RowsAffected()
					if err != nil {
						resWork.err = err.Error()
						responseQueue <- resWork
						continue
					} else if a == 0 {
						resWork.err = "Khong tim thay thue bao yeu cau hoac khong co gi thay doi thong tin thue bao\n"
						responseQueue <- resWork
						continue
					}
				}
			} else if work.CMD == 3 {
				stmt, err := db.Prepare("DELETE FROM infor_customer WHERE MSISDN = ?")
				if err != nil {
					resWork.err = err.Error()
					responseQueue <- resWork
					continue
				} else {
					res, err := stmt.Exec(getCus.MSISDN)
					a, _ := res.RowsAffected()
					if err != nil {
						resWork.err = err.Error()
						responseQueue <- resWork
						continue
					} else if a == 0 {
						resWork.err = "Khong tim thay thue bao yeu cau\n"
						responseQueue <- resWork
						continue
					}
				}
			} else {
				resWork.err = "Khong dung dinh dang\n"
				responseQueue <- resWork
				continue
			}
			fmt.Println("b")
			resWork.err = "0"
			responseQueue <- resWork
			continue
		} else {
			resWork.err = err.Error()
			responseQueue <- resWork
			fmt.Println(err)
			continue
		}
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
