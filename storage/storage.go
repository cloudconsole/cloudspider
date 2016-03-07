package storage

import (
	"gopkg.in/mgo.v2"
	"github.com/spf13/viper"
	log "github.com/Sirupsen/logrus"
	"strings"
	"sync"
)

var MongoSession *mgo.Session

func DBConnection() *mgo.Database {
	if MongoSession == nil {
		uri := viper.GetString("mongodb.uri")
		if uri == "" {
			log.Fatal("No connection uri for MongoDB provided")
		}

		var err error
		MongoSession, err = mgo.Dial(uri)
		if MongoSession == nil || err != nil {
			log.Fatal("Can't connect to mongo, go error %v\n", err)
		}

		MongoSession.SetSafe(&mgo.Safe{})
	}

	return MongoSession.DB(viper.GetString("mongodb.dbname"))
}

// Write many docs in to mongodb
func InsertMany(Docs []interface{}, collName string, wg *sync.WaitGroup) {
	coll := DBConnection().C(collName)

	for _, doc := range Docs {
		if err := coll.Insert(doc); err != nil {
			if !strings.Contains(err.Error(), "E11000") {
				log.Fatal("Database error. Err: %v", err)
			}
		}
		wg.Done()  // say write is done
	}
}

// Checks if the document exists with the given ID
func DocExists(Id string, collName string) bool {
	coll := DBConnection().C(collName)

	query := coll.FindId(Id)
	count, err := query.Count()
	if err != nil {
		log.Fatal("Database error. Err: %v", err)
	}

	if count >= 1 {
		return true
	} else {
		return false
	}
}

// Ensure index exits
func EnsureMachinesIndex() error {
	coll := DBConnection().C("machines")

	return coll.EnsureIndex(
		mgo.Index{
			Key: []string{
				"$text:ssh_key_name",
				"$text:type",
				"$text:public_dns",
				"$text:private_dns",
				"$text:tags",
				"$text:security_group",
			},
			Name: "machines_index",
			Background: true,
		})
}
