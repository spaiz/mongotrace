package app

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/spaiz/mongotrace/app/formatters"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

var operationBgColor = color.HEX("#B63B2E", true)
var timestampTextColor = color.Black.Render
var timestampBgColor = color.BgWhite.Render

type Formatter interface {
	Format(row *LogResult) (string, error)
}

func NewDefaultOplogFormatter() *OplogFormatter {
	return &OplogFormatter{}
}

type OplogFormatter struct{}

func (receiver *OplogFormatter) Format(logResults *LogResult) (string, error) {
	formatted := ""
	collection := logResults.Namespace
	operation := operationBgColor.Sprintf(strings.ToUpper(" "+logResults.Op) + " ")
	rt := fmt.Sprintf("%s", timestampBgColor(timestampTextColor(fmt.Sprintf(" %s ", logResults.ExecutionTime))))

	var formatter app.CommandFormatter

	var cmd bson.M
	switch logResults.Op {
	case "command", "query":
		cmd = logResults.Query
		if cmd == nil {
			cmd = logResults.Command
		}

		_, ok := cmd["$readPreference"]
		if ok {
			return "", nil
		}

		var err error
		formatter, err = app.GetCommandFormatter(cmd, logResults.Op)
		if err != nil {
			result, err := json.Marshal(logResults)
			if err != nil {
				return "", err
			}

			formatted = fmt.Sprintf("Missing formatter: %s: %s\n", operation, string(result))
		}

	case "getmore":
		cmd = logResults.OriginatingCommand
		var err error
		formatter, err = app.GetCommandFormatter(cmd, logResults.Op)
		if err != nil {
			result, err := json.Marshal(logResults)
			if err != nil {
				return "", err
			}

			formatted = fmt.Sprintf("Missing formatter: %s: %s\n", operation, string(result))
		}

	default:
		result, err := json.Marshal(logResults)
		if err != nil {
			return "", err
		}

		fmt.Printf("Not supported operation received %s: %s\n", operation, string(result))
		return formatted, nil
	}

	if formatter == nil {
		return formatted, nil
	}

	str, err := formatter(cmd, collection)
	if err != nil {
		return "", nil
	}

	formatted = fmt.Sprintf("%s%s %s\n", operation, rt, str)
	return formatted, nil
}
