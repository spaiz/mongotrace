package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func AggregateFormat(cmdDoc bson.M, collection string) (string, error) {
	pipeline := cmdDoc["pipeline"].(bson.A)
	options := bson.M{}

	optionsFields :=[]string{
		"explain",
		"allowDiskUse",
		"cursor",
		"maxTimeMS",
		"bypassDocumentValidation",
		"readConcern",
		"collation",
		"hint",
		"comment",
		"writeConcern",
	}

	for _, field := range optionsFields {
		val, ok := cmdDoc[field]
		if ok {
			options[field] = val
		}
	}

	pipelineParam, err := json.Marshal(pipeline)
	if err != nil {
		return "", err
	}

	optionsParam, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	queryFormatted := fmt.Sprintf("%s.aggregate(%s, %s)", collection, pipelineParam, optionsParam)
	return queryFormatted, nil
}