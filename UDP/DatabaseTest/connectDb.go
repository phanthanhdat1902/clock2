package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type customer struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	MSISDN   string        `bson:"MSISDN"`
	IMSI     string        `bson:"IMSI"`
	Name     string        `bson:"name"`
	CMT      string        `bson:"CMT"`
	Birthday string        `bson:"birthday"`
}

func main() {
	collection := OpenDatabase()
	var john customer
	john.MSISDN = "84981064189"
	john.IMSI = "12345642"
	john.Name = "DAT"
	john.Birthday = "19021998"
	john.CMT = "123456780"
	r, e := collection.InsertOne(context.Background(), john)
	fmt.Println(r)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("sucessfull")
	}
}
func OpenDatabase() *mongo.Collection {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	// Connect to MongoDB
	client, e := mongo.Connect(context.TODO(), clientOptions)
	if e != nil {
		fmt.Println(e)
	}
	// Check the connection
	e = client.Ping(context.TODO(), nil)
	if e != nil {
		fmt.Println(e)
	}
	collection := client.Database("my_exam").Collection("customer")
	return collection
}
