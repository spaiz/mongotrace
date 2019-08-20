package app

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type FormatterNotFoundError struct {
	Action string
}

func (r FormatterNotFoundError) Error() string {
	return fmt.Sprintf("formatter for the %s not found", r.Action)
}

type CommandFormatter func(cmdDoc bson.M, collection string) (string, error)

type FormattersRegistry map[string]CommandFormatter

var (
	formatters = FormattersRegistry{
		"find": FindFormat,
		"distinct": DistinctFormat,
		"count": CountFormat,
		"aggregate": AggregateFormat,
		"remove": RemoveFormat,
		"findandmodify": FindAndModifyFormat,
		"listIndexes": ListIndexesFormat,
		"group": GroupFormat,
		"update": UpdateFormat,
		"insert": InsertFormat,
		"getmore": GetmoreFormat,
	}
)

func GetCommandFormatter(cmd bson.M, action string) (CommandFormatter, error) {
	if action == "getmore" {
		formatter, ok := formatters[action]
		if ok {
			return formatter, nil
		}
	}

	for key, _ := range formatters {
		_, ok := cmd[key]
		if ok {
			formatter, ok := formatters[key]
			if ok {
				return formatter, nil
			}
		}
	}

	return nil, &FormatterNotFoundError{
		Action: action,
	}
}
