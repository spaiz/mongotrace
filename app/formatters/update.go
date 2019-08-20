package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateFormat(cmdDoc bson.M, collection string) (string, error) {
	var qJSON []byte
	q, ok := cmdDoc["q"]
	if ok {
		if q != nil {
			var err error
			qJSON, err = json.Marshal(q)
			if err != nil {
				return "", err
			}
		}
	}

	var uJSON []byte
	u, ok := cmdDoc["u"]
	if ok {
		if u != nil {
			var err error
			uJSON, err = json.Marshal(u)
			if err != nil {
				return "", err
			}
		}
	}

	options := map[string]interface{}{}

	multi, ok := cmdDoc["multi"]
	if ok {
		options["multi"] = multi
	}

	upsert, ok := cmdDoc["upsert"]
	if ok {
		options["upsert"] = upsert
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.update(%s, %s, %s)", collection, qJSON, uJSON, optionsJSON), nil
}