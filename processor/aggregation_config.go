package processor

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	AGGREGATED_CONFIG_PATH = "./aggregation_config"
)

var (
	loadConfig              = loadAggregatedConfig()
	rootAggregationFunction = map[string]map[string]*AggregatedTask{}
)

func loadAggregatedConfig() bool {
	files, _ := ioutil.ReadDir(AGGREGATED_CONFIG_PATH)

	for _, file := range files {
		jsonFile, _ := os.Open(AGGREGATED_CONFIG_PATH + "/" + file.Name())
		byteValue, _ := ioutil.ReadAll(jsonFile)
		config := AggregatedConfig{}
		json.Unmarshal(byteValue, &config)
		newAggregatedTask(&config)
	}

	return true
}

func newAggregatedTask(config *AggregatedConfig) {
	task := AggregatedTask{
		ID:           config.AggregatedID,
		GroupByKeyID: config.GroupByKeyID,
	}

	for _, f := range config.AggregatedFunction {
		aggFunc := AggregatedFunction{
			PropertyName: f.PropertyName,
			FieldName:    f.FieldName,
		}

		if function := getFunctionImpl(f.Function); function != nil {
			aggFunc.Function = function
			task.Functions = append(task.Functions, aggFunc)
		}
	}

	registerTask(config.DomainName, &task)
}

func getFunctionImpl(funcName string) Function {
	if funcName == "sum" {
		return &SumFunction{}
	} else if funcName == "min" {
		return &MinFunction{}
	} else if funcName == "max" {
		return &MaxFunction{}
	} else if funcName == "count" {
		return &CountFunction{}
	} else if funcName == "last" {
		return &LastFunction{}
	}

	return nil
}

func registerTask(domain string, task *AggregatedTask) {
	aggregation, ok := rootAggregationFunction[domain]

	if !ok {
		aggregation = map[string]*AggregatedTask{}
		rootAggregationFunction[domain] = aggregation
	}

	aggregation[task.ID] = task
}
