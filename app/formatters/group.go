package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func GroupFormat(cmdDoc bson.M, collection string) (string, error) {
	group := cmdDoc["group"].(bson.M)

	keyParam, err := json.Marshal(group["key"])
	if err != nil {
		return "", err
	}

	initialParam, err := json.Marshal(group["initial"])
	if err != nil {
		return "", err
	}

	condParam, err := json.Marshal(group["cond"])
	if err != nil {
		return "", err
	}

	collationParam, err := json.Marshal(group["collation"])
	if err != nil {
		return "", err
	}

	reduceParam := ""
	reduce, ok := group["$reduce"]
	if ok {
		reduceParam = fmt.Sprintf("%s", reduce)
	}

	keyfParam := ""
	keyf, ok := group["keyf"]
	if ok {
		keyfParam = fmt.Sprintf("%s", keyf)
	}

	finalizeParam := ""
	finalize, ok := group["finalize"]
	if ok {
		finalizeParam = fmt.Sprintf("%s", finalize)
	}

	groupParams := `{ key: %s, reduce: %s, initial: %s, keyf: %s, cond: %s, finalize: %s, collation: %s }`
	groupParams = fmt.Sprintf(groupParams, keyParam, reduceParam, initialParam, keyfParam, condParam, finalizeParam, collationParam)
	queryFormatted := fmt.Sprintf("%s.group(%s)", collection, groupParams)
	return queryFormatted, nil
}