package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func GetmoreFormat(cmdDoc bson.M, collection string) (string, error) {
	queryFormatted := ""

	var documentsJSON []byte
	documents, documentsExists := cmdDoc["documents"]
	if documentsExists {
		if documents != nil {
			var err error
			documentsJSON, err = json.Marshal(documents)
			if err != nil {
				return "", err
			}
		}
	}

	options := map[string]interface{}{}
	writeConcern, ok := cmdDoc["writeConcern"]
	if ok {
		options["writeConcern"] = writeConcern
	}

	ordered, ok := cmdDoc["ordered"]
	if ok {
		options["ordered"] = ordered
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	queryFormatted = fmt.Sprintf("%s.insert(%s, %s)", collection, string(documentsJSON), optionsJSON)
	return queryFormatted, nil
}
