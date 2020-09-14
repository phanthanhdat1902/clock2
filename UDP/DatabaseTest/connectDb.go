package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)
import (
	_ "gopkg.in/mgo.v2"
)

type customer struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	MSISDN   string        `bson:"MSISDN"`
	IMSI     string        `bson:"IMSI"`
	Name     string        `bson:"name"`
	CMT      string        `bson:"CMT"`
	Birthday string        `bson:"birthday"`
}

const (
	MongoDBHosts = "mongodb://127.0.0.1:27017"
	AuthDatabase = "my_exam"
)

func main() {
	mongoSession := OpenDatabase()
	mongoSession.SetMode(mgo.Monotonic, true)
	var john customer
	john.MSISDN = "84981064189"
	john.IMSI = "12345642"
	john.Name = "DAT"
	john.Birthday = "19021998"
	john.CMT = "123456780"
	insert(mongoSession, john)
}
func OpenDatabase() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
	}
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	return mongoSession
}
func insert(mongoSession *mgo.Session, john customer) {
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("my_exam").C("customer")
	err := collection.Insert(john)
	fmt.Println(err)
}
