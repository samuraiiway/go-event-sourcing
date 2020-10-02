package repository

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	PREFIX_EVENT_DB_PATH = "./db/event/"
)

var (
	rootEventDB = map[string]*leveldb.DB{}
)

func getEventDB(domain string) *leveldb.DB {
	db, ok := rootEventDB[domain]
	if !ok {
		db, _ = leveldb.OpenFile(PREFIX_EVENT_DB_PATH+domain, nil)
		rootEventDB[domain] = db
	}

	return db
}

func SaveEvent(domain string, key string, event map[string]interface{}) {
	db := getEventDB(domain)
	value, _ := json.Marshal(event)
	db.Put([]byte(key), value, nil)
}
