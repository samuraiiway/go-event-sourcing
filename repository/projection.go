package repository

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	PREFIX_PROJECTION_DB_PATH = "./db/projection/"
	PROJECTION_CACHE_SIZE     = 100000
)

var (
	rootProjectionDB    = map[string]*leveldb.DB{}
	rootProjectionCache = map[string]map[string]map[string]interface{}{}
)

func getDB(domain string) *leveldb.DB {
	db, ok := rootProjectionDB[domain]
	if !ok {
		db, _ = leveldb.OpenFile(PREFIX_PROJECTION_DB_PATH+domain, nil)
		rootProjectionDB[domain] = db
	}

	return db
}

func getCache(domain string) map[string]map[string]interface{} {
	cache, ok := rootProjectionCache[domain]

	if !ok {
		cache = make(map[string]map[string]interface{}, PROJECTION_CACHE_SIZE)
		rootProjectionCache[domain] = cache
	}

	return cache
}

func cacheProjection(domain string, key string, projection map[string]interface{}) {
	cache := getCache(domain)
	if len(cache) == PROJECTION_CACHE_SIZE {
		for key, _ := range cache {
			delete(cache, key)
			break
		}
	}
	cache[key] = projection
}

func persistProjection(domain string, key string, projection map[string]interface{}) {
	db := getDB(domain)
	value, _ := json.Marshal(projection)
	db.Put([]byte(key), value, nil)
}

func GetProjection(domain string, key string) map[string]interface{} {
	cache := getCache(domain)
	projection, ok := cache[key]
	if !ok {
		db := getDB(domain)
		bytes, err := db.Get([]byte(key), nil)
		if err == nil {
			json.Unmarshal(bytes, &projection)
		}
	}

	if projection == nil {
		projection = map[string]interface{}{}
	}

	return projection
}

func SaveProjection(domain string, key string, projection map[string]interface{}) {
	cacheProjection(domain, key, projection)
	persistProjection(domain, key, projection)
}
