package app

import (
	"encoding/json"
	"fmt"
	"github.com/spaiz/mongotrace/config"
	"github.com/spaiz/mongotrace/db"
	"github.com/tidwall/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sync"
	"time"
)

type LogResult struct {
	AppName            string         `bson:"appName"`
	Command            bson.M         `bson:"command"`
	Client             string         `bson:"client"`
	ExecutionTime      *time.Duration `bson:"millis"`
	Namespace          string         `bson:"ns"`
	Op                 string         `bson:"op"`
	Query              bson.M         `bson:"query"`
	Time               time.Time      `bson:"ts"`
	Update             bson.M         `bson:"update"`
	Upsert             bool           `bson:"upsert"`
	User               string         `bson:"user"`
	OriginatingCommand bson.M `bson:"originatingCommand"`
}

func NewLogReportCommand(conf *config.Configuration, options *Options, dbs []*db.MongoDB, formatter Formatter) *LogReportCommand {
	return &LogReportCommand{
		conf:    conf,
		options: options,
		dbs:     dbs,
		formatter: formatter,
	}
}

type LogReportCommand struct {
	conf    *config.Configuration
	options *Options
	dbs     []*db.MongoDB
	formatter Formatter
}

func (r *LogReportCommand) Execute() error {
	hostsNum := len(r.dbs)
	var wg sync.WaitGroup
	wg.Add(hostsNum)

	out := make(chan *LogItem, 10000)

	go r.receive(out)

	for _, dbconn := range r.dbs {
		go func(dbd *db.MongoDB) {
			defer wg.Done()
			reporter := NewLogReporter(dbd, out)
			fmt.Printf("Run reporter of %v\n", dbd.GetHost())
			err := reporter.Run()
			if err != nil {
				log.Fatalf("Failed to run reporter; %v\n", err.Error())
			}
		}(dbconn)
	}

	wg.Wait()
	return nil
}

/*
	Receives oplog documents, formats it and prints to stdout
 */
func (r *LogReportCommand) receive(out chan *LogItem) {
	for item := range out {
		logResult := LogResult{}
		err := bson.Unmarshal(item.RawData, &logResult)

		if err != nil {
			log.Printf("Error when parsing mongodb query: %s", err.Error())
			continue
		}

		if r.options.Raw || r.options.Debug {
			rawResults := map[string]interface{}{}
			err := bson.Unmarshal(item.RawData, &rawResults)

			result, err := json.Marshal(rawResults)
			if err != nil {
				log.Printf("Error: %s", err.Error())
				continue
			}

			if r.options.Indent {
				result = pretty.PrettyOptions(result, &pretty.Options{
					Width:    0,
					Prefix:   "",
					Indent:   "  ",
					SortKeys: false,
				})
			}

			if r.options.Color {
				result = pretty.Color(result, nil)
			}

			if r.options.Indent {
				fmt.Printf("%s", string(result))
			} else {
				fmt.Printf("%s\n", string(result))
			}

			if !r.options.Debug {
				continue
			}
		}

		str, err := r.formatter.Format(&logResult)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			continue
		}

		if str != "" {
			fmt.Printf("%s", str)
		}
	}
}
