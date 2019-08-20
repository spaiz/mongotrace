package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func DefaultFormat(cmdDoc bson.M, collection string) (string, error) {
	params := &distinctFormatParams{
		Field: cmdDoc["key"],
		Query: cmdDoc["query"],
		Options: bson.M{
			"collation": cmdDoc["collation"],
		},
	}

	queryParam, err := json.Marshal(params.Query)
	if err != nil {
		return "", err
	}

	optionsParam, err := json.Marshal(params.Options)
	if err != nil {
		return "", err
	}

	queryFormatted := fmt.Sprintf("%s.distinct(\"%s\", %s, %s)", collection, params.Field, queryParam, optionsParam)
	return queryFormatted, nil
}