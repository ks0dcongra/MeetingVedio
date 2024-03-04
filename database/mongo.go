package database

import (
	"context"
	"github.com/qiniu/qmgo"
)

var QmgoConnection *qmgo.QmgoClient
var err error

func MongoDBinit() {
	ctx := context.TODO()
	QmgoConnection, err = qmgo.Open(ctx, &qmgo.Config{
		Uri: "mongodb://localhost:27017/",
		// Uri:      "mongodb://root:root@db:27017/",
		// Uri:      "mongodb://root:root@172.31.12.7:27017/",
		// Uri:      "mongodb://root:root@172.31.18.14:27017/",
		// Uri:      "mongodb://root:root@172.31.3.255:27017/",
		Database: "meetingVideoHelper",
		Coll:     "videoHelper",
	})
	if err != nil {
		panic(err)
	}
}
