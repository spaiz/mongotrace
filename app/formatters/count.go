package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func CountFormat(cmdDoc bson.M, collection string) (string, error) {
	options := bson.M{}
	optionsFields :=[]string{
		"limit",
		"skip",
		"hint",
		"maxTimeMS",
		"readConcern",
	}

	for _, field := range optionsFields {
		val, ok := cmdDoc[field]
		if ok {
			options[field] = val
		}
	}

	query := cmdDoc["query"]
	queryParam, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	optionsParam, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	queryFormatted := fmt.Sprintf("%s.count(%s, %s)", collection, queryParam, optionsParam)
	return queryFormatted, nil
}