package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type insertParams struct {
	WriteConcern interface{} `json:"writeConcern,omitempty"`
	Ordered      interface{} `json:"sort,omitempty"`
}

func InsertFormat(cmdDoc bson.M, collection string) (string, error) {
	params := &insertParams{
		WriteConcern: cmdDoc["writeConcern"],
		Ordered:      cmdDoc["ordered"],
	}

	queryFormatted := ""

	var documentsJSON []byte
	documents, ok := cmdDoc["documents"]
	if ok {
		if documents != nil {
			var err error
			documentsJSON, err = json.Marshal(documents)
			if err != nil {
				return "", err
			}
		}
	}

	optionsJSON, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	queryFormatted = fmt.Sprintf("%s.insert(%s, %s)", collection, string(documentsJSON), optionsJSON)
	return queryFormatted, nil
}

