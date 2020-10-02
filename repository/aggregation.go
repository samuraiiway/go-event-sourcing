package repository

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	PREFIX_AGGREGATION_DB_PATH = "./db/aggregation/"
	AGGREGATION_CACHE_SIZE     = 100000
)

var (
	rootAggregationDB    = map[string]*leveldb.DB{}
	rootAggregationCache = map[string]map[string]map[string]interface{}{}
)

func getAggregationDB(domain string) *leveldb.DB {
	db, ok := rootAggregationDB[domain]
	if !ok {
		db, _ = leveldb.OpenFile(PREFIX_AGGREGATION_DB_PATH+domain, nil)
		rootAggregationDB[domain] = db
	}

	return db
}

func getAggregationCache(domain string) map[string]map[string]interface{} {
	cache, ok := rootAggregationCache[domain]

	if !ok {
		cache = make(map[string]map[string]interface{}, AGGREGATION_CACHE_SIZE)
		rootAggregationCache[domain] = cache
	}

	return cache
}

func cacheAggregation(domain string, key string, aggregation map[string]interface{}) {
	cache := getAggregationCache(domain)
	if len(cache) == AGGREGATION_CACHE_SIZE {
		for key, _ := range cache {
			delete(cache, key)
			break
		}
	}
	cache[key] = aggregation
}

func persistAggregation(domain string, key string, aggregation map[string]interface{}) {
	db := getAggregationDB(domain)
	value, _ := json.Marshal(aggregation)
	db.Put([]byte(key), value, nil)
}

func GetAggregation(domain string, key string) map[string]interface{} {
	cache := getAggregationCache(domain)
	aggregation, ok := cache[key]
	if !ok {
		db := getAggregationDB(domain)
		bytes, err := db.Get([]byte(key), nil)
		if err == nil {
			json.Unmarshal(bytes, &aggregation)
		}
	}

	if aggregation == nil {
		aggregation = map[string]interface{}{}
	}

	return aggregation
}

func SaveAggregation(domain string, key string, aggregation map[string]interface{}) {
	cacheAggregation(domain, key, aggregation)
	persistAggregation(domain, key, aggregation)
}
