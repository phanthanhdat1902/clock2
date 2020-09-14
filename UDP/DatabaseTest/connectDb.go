package main

import (
	"fmt"
)
import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
	session, _ := mgo.Dial("127.0.0.1:27017")
	return session
}
func insert(mongoSession *mgo.Session, john customer) {
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("my_exam").C("customer")
	err := collection.Insert(john)
	fmt.Println(err)
}
