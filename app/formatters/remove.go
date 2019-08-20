package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func RemoveFormat(cmdDoc bson.M, collection string) (string, error) {
	q := cmdDoc["q"].(bson.M)
	options := bson.M{}
	limit := cmdDoc["limit"]

	optionsFields :=[]string{
		"writeConcern",
		"collation",
	}

	for _, field := range optionsFields {
		val, ok := cmdDoc[field]
		if ok {
			options[field] = val
		}
	}

	qParam, err := json.Marshal(q)
	if err != nil {
		return "", err
	}

	optionsParam, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	queryFormatted := fmt.Sprintf("%s.remove(%s, %s).limit(%v)", collection, qParam, optionsParam, limit)
	return queryFormatted, nil
}
