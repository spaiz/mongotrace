package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func FindFormat(cmdDoc bson.M, collection string) (string, error) {
	queryFormatted := ""

	var filterJSON []byte
	filter, ok := cmdDoc["filter"]
	if ok {
		var err error
		filterJSON, err = json.Marshal(filter)
		if err != nil {
			return "", err
		}
	}

	if string(filterJSON) == "" {
		filterJSON = []byte("{}")
	}

	find, findExists := cmdDoc["find"]
	if findExists {
		queryFormatted = fmt.Sprintf("%s.find(%s)", find, filterJSON)
	}

	projection, ok := cmdDoc["projection"]
	if ok {
		projectionJSON, err := json.Marshal(projection)
		if err != nil {
			return "", err
		}

		if string(projectionJSON) != "{}" {
			queryFormatted = fmt.Sprintf("%s.projection(%s)", queryFormatted, projectionJSON)
		}
	}

	sort, ok := cmdDoc["sort"]
	if ok {
		if sort != nil {
			sortJSON, err := json.Marshal(sort)
			if err != nil {
				return "", err
			}

			queryFormatted = fmt.Sprintf("%s.sort(%s)", queryFormatted, sortJSON)
		}
	}

	limit, ok := cmdDoc["limit"]
	if ok {
		if limit != nil {
			limitJSON, err := json.Marshal(limit)
			if err != nil {
				return "", err
			}

			queryFormatted = fmt.Sprintf("%s.limit(%s)", queryFormatted, limitJSON)
		}
	}

	return queryFormatted, nil
}