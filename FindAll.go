package main

import (
	"context"
	"fmt"
	"MeetingVideoHelper/database"
	"MeetingVideoHelper/models"

	"go.mongodb.org/mongo-driver/bson"
)

func main5() {
	user := []models.MongoUser{}
	database.QmgoConnection.Find(context.TODO(), bson.M{"name":"Tom"}).All(&user)
	for _, user := range user {
		fmt.Printf("%+v",user)
	}
}