package storage

import (
	"github.com/cloudconsole/cloudspider/log"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

var MongoSession *mgo.Session

func DBConnection() *mgo.Database {
	if MongoSession == nil {
		uri := viper.GetString("mongodb.uri")
		if uri == "" {
			log.Fatal(map[string]interface{}{},
				"No connection uri for MongoDB provided")
		}

		var err error
		MongoSession, err = mgo.Dial(uri)
		if MongoSession == nil || err != nil {
			log.Fatal(map[string]interface{}{
				"DBError": err,
			}, "Can't connect to mongo")
		}

		MongoSession.SetSafe(&mgo.Safe{})
	}

	return MongoSession.DB(viper.GetString("mongodb.dbname"))
}
