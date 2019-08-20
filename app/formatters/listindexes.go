package app

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func ListIndexesFormat(_ bson.M, collection string) (string, error) {
	queryFormatted := fmt.Sprintf("%s.getIndexes()", collection)
	return queryFormatted, nil
}
