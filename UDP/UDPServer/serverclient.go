package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"runtime"
	"strconv"
)
import _ "github.com/go-sql-driver/mysql"

const (
	DB_HOST = "tcp(127.0.0.1:3306)"
	DB_NAME = /*name database*/ "my_exam"
	DB_USER = /*"user"*/ "ptd"
	DB_PASS = /*"pass"*/ "anh123asd"
)

var number int = 0

type customer struct {
	MSISDN   string
	IMSI     string
	name     string
	CMT      string
	birthday string
}

var done = make(chan bool)

/*open*/
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME))
	if err != nil {
		fmt.Println(err)
		return db, err
	}
	return db, nil
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var saddr net.UDPAddr
	saddr.IP = net.ParseIP("127.0.0.1")
	saddr.Port = 8888
	connection, _ := net.ListenUDP("udp", &saddr)
	for i := 0; i < runtime.NumCPU(); i++ {
		go handleConnect(connection)
	}
	fmt.Println("server running")
	<-done
	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	msg, _ := reader.ReadString('\n')
	//	if strings.Compare(msg, "close server\n") == 0 {
	//		break
	//	}
	//}
}

func handleConnect(connection *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		memset(buffer)
		_, client, _ := connection.ReadFromUDP(buffer)
		go handleMsg(string(buffer), connection, client)
		number++
		fmt.Println(number)
	}
}

func handleMsg(msg string, connection *net.UDPConn, client *net.UDPAddr) {
	var cus customer
	//get CMD_MSISDN
	CMD_MSISDN := msg[:6]
	var result error
	chanCus := make(chan customer)
	chanCMD := make(chan int)
	chanErr := make(chan error)
	go workingDatabase(chanCus, chanCMD, chanErr)
	CMD_MSISDN = hex.EncodeToString([]byte(CMD_MSISDN))
	CMD, _ := strconv.Atoi(string(CMD_MSISDN[0]))
	MSISDN := CMD_MSISDN[1:]
	cus.MSISDN = MSISDN
	fmt.Println(MSISDN)
	//get IMSI
	lenIMSI, _ := strconv.Atoi(hex.EncodeToString([]byte(string(msg[6]))))
	lenIMSI -= 10
	index := 7 + (lenIMSI+1)/2
	IMSI := msg[7:index]
	IMSI = hex.EncodeToString([]byte(IMSI))
	IMSI = IMSI[:lenIMSI]
	cus.IMSI = IMSI
	if CMD == 3 {
		//func remove
		chanCus <- cus
		chanCMD <- 3
		result = <-chanErr
		fmt.Println("handleMsg: ", chanErr)
		if result != nil {
			connection.WriteToUDP([]byte("400 Err\n"), client)
		} else {
			connection.WriteToUDP([]byte("200 OK\n"), client)
		}
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
		chanCus <- cus
		chanCMD <- CMD
		result = <-chanErr
		if result != nil {
			connection.WriteToUDP([]byte("400 Err\n"), client)
		} else {
			connection.WriteToUDP([]byte("200 OK\n"), client)
		}
	}
}

func workingDatabase(chanCus chan customer, chanCMD chan int, chanErr chan error) {
	getCus := <-chanCus
	getCMD := <-chanCMD
	db, err := OpenDB()
	if db != nil {
		if getCMD == 1 {
			stmt, err := db.Prepare("insert into infor_customer values (?,?,?,?,?)")
			if err != nil {
				chanErr <- err
				fmt.Println("workingDtabase db.Prepare", err)
			} else {
				_, err := stmt.Exec(getCus.MSISDN, getCus.IMSI, getCus.name, getCus.CMT, getCus.birthday)
				db.Close()
				if err != nil {
					chanErr <- err
				}
			}
		} else if getCMD == 2 {
			stmt, err := db.Prepare("update infor_customer set full_name=?,CMND=?,birthday=? where MSISDN=?")
			if err != nil {
				chanErr <- err
			} else {
				_, err := stmt.Exec(getCus.name, getCus.CMT, getCus.birthday, getCus.MSISDN)
				db.Close()
				if err != nil {
					chanErr <- err
				}
			}
		} else if getCMD == 3 {
			stmt, err := db.Prepare("DELETE FROM infor_customer WHERE MSISDN = ?")
			if err != nil {
				chanErr <- err
			} else {
				_, err := stmt.Exec(getCus.MSISDN)
				db.Close()
				if err != nil {
					chanErr <- err
				}
			}
		}
		chanErr <- nil
		return
	} else {
		chanErr <- err
		fmt.Println(err)
		return
	}
}

//func workingDatabase(cus chan customer,cmd chan int,result chan bool)  {
//	db:=OpenDB()
//	defer db.Close()
//	if db!=nil{
//		getCus:=<-cus
//		getCMD:=<-cmd
//		if getCMD==1{
//			stmt, err := db.Prepare("insert into infor_customer values (?,?,?,?,?)")
//			if err!=nil{
//				result<-false
//			}else{
//				_,err:=stmt.Exec(getCus.MSISDN,getCus.IMSI,getCus.name,getCus.CMT,getCus.birthday)
//				if err!=nil{
//					result<-false
//				}
//			}
//		}else if getCMD==2{
//			stmt, err := db.Prepare("update infor_customer set full_name=?,CMND=?,birthday=? where MSISDN=?")
//			if err!=nil{
//				result<-false
//			}else{
//				_,err:=stmt.Exec(getCus.name,getCus.CMT,getCus.birthday,getCus.MSISDN)
//				if err!=nil{
//					result<-false
//				}
//			}
//		}else if getCMD==3{
//			stmt, err := db.Prepare("DELETE FROM infor_customer WHERE MSISDN = ?")
//			if err!=nil{
//				result<-false
//			}else{
//				_,err:=stmt.Exec(getCus.MSISDN)
//				if err!=nil{
//					result<-false
//				}
//			}
//		}
//		result<-true
//	}else{
//		result<-false
//	}
//}
//func create(cus customer) error {
//	db:=OpenDB()
//	if db!=nil{
//		stmt, err := db.Prepare("insert into infor_customer values (?,?,?,?,?)")
//		if err!=nil{
//			return err
//		}else{
//			_,err:=stmt.Exec(cus.MSISDN,cus.IMSI,cus.name,cus.CMT,cus.birthday)
//			if err!=nil{
//				return err
//			}
//		}
//	}else{
//		err := db.Ping()
//		return err
//	}
//	defer db.Close()
//	return nil
//}
////func read()  {
////
////}
//func update(cus customer)error  {
//	db:=OpenDB()
//	if db!=nil{
//		stmt, err := db.Prepare("update infor_customer set full_name=?,CMND=?,birthday=? where MSISDN=?")
//		if err!=nil{
//			return err
//		}else{
//			_,err:=stmt.Exec(cus.name,cus.CMT,cus.birthday,cus.MSISDN)
//			if err!=nil{
//				return err
//			}
//		}
//	}else{
//		err := db.Ping()
//		return err
//	}
//	defer db.Close()
//	return nil
//}
//func delete(cus customer) error  {
//	db:=OpenDB()
//	if db!=nil{
//		stmt, err := db.Prepare("DELETE FROM infor_customer WHERE MSISDN = ?")
//		if err!=nil{
//			return err
//		}else{
//			_,err:=stmt.Exec(cus.MSISDN)
//			if err!=nil{
//				return err
//			}
//		}
//	}else{
//		err := db.Ping()
//		return err
//	}
//	defer db.Close()
//	return nil
//}

func memset(des []byte) {
	for i := range des {
		des[i] = byte(0)
	}
}
