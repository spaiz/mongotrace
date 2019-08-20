package app

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type findAndModifyFormatParams struct {
	Query                    interface{} `json:"query,omitempty"`
	Sort                     interface{} `json:"sort,omitempty"`
	Remove                   interface{} `json:"remove,omitempty"`
	Update                   interface{} `json:"update,omitempty"`
	New                      interface{} `json:"new,omitempty"`
	Fields                   interface{} `json:"fields,omitempty"`
	Upsert                   interface{} `json:"upsert,omitempty"`
	BypassDocumentValidation interface{} `json:"bypassDocumentValidation,omitempty"`
	WriteConcern             interface{} `json:"writeConcern,omitempty"`
	Collation                interface{} `json:"collation,omitempty"`
	ArrayFilters             interface{} `json:"arrayFilters,omitempty"`
}

func FindAndModifyFormat(cmdDoc bson.M, collection string) (string, error) {
	params := &findAndModifyFormatParams{
		Query:                    cmdDoc["query"],
		Sort:                     cmdDoc["sort"],
		Remove:                   cmdDoc["remove"],
		Update:                   cmdDoc["update"],
		New:                      cmdDoc["new"],
		Fields:                   cmdDoc["fields"],
		Upsert:                   cmdDoc["upsert"],
		BypassDocumentValidation: cmdDoc["bypassDocumentValidation"],
		WriteConcern:             cmdDoc["writeConcern"],
		Collation:                cmdDoc["collation"],
		ArrayFilters:             cmdDoc["arrayFilters"],
	}

	funcParams, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	queryFormatted := fmt.Sprintf("%s.findAndModify(%s)", collection, funcParams)
	return queryFormatted, nil
}