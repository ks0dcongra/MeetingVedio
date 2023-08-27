package database

import (
	"context"
	"github.com/qiniu/qmgo"
)

var QmgoConnection *qmgo.QmgoClient
var err error

func init() {
	ctx := context.TODO()
	QmgoConnection, err = qmgo.Open(ctx, &qmgo.Config{
		Uri:      "mongodb://localhost:27017/",
		Database: "meetingVideoHelper",
		Coll:     "videoHelper",
	})
	if err != nil {
		panic(err)
	}
}
