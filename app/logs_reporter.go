package app

import (
	"context"
	"fmt"
	"github.com/spaiz/mongotrace/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogItem struct {
	RawData []byte
	Src *LogsReporter
}

type Profile struct {
	Was    int `bson:"was"`
	Slowms int `bson:"slowms"`
	Ok     int `bson:"ok"`
	DbName string
	Host   string
}

/*
	Defined LogsReporter struct
 */
type LogsReporter struct {
	DB            *db.MongoDB
	onDataHandler func(data []byte)
	fields        []string
	dbName        string
	nsIgnored     []string
	out chan *LogItem
}

/*
	Returns New LogsReporter instance.
	LogsReporter used to fetch logs from a single database
 */
func NewLogReporter(DB *db.MongoDB, out chan *LogItem) *LogsReporter {
	fields := []string{"ts", "millis", "op", "ns", "query", "updateobj", "command", "ninserted", "ndeleted", "nMatched", "nreturned"}

	var nsIgnored []string
	nsIgnored = append(nsIgnored, fmt.Sprintf("%s.system.profile", DB.Name()))

	return &LogsReporter{
		DB:        DB,
		fields:    fields,
		dbName:    DB.Name(),
		nsIgnored: nsIgnored,
		out: out,
	}
}

// Run used to start the logs pulling
func (r *LogsReporter) Run() error {
	for {
		err := r.findLogs()
		if err != nil {
			return err
		}
	}
}

// Opens cursor to oplog collection, and send received documents to channel for post processing
func (r *LogsReporter) findLogs() error {
	query := bson.D{}

	el := bson.E{
		Key: "ns",
		Value: bson.M{
			"$nin": r.nsIgnored,
		},
	}

	query = append(query, el)

	findOptions := &options.FindOptions{}
	findOptions.SetCursorType(options.TailableAwait)
	findOptions.SetMaxTime(10 * time.Second)
	findOptions.SetBatchSize(100)

	projection := bson.D{}

	for _, field := range r.fields {
		el := bson.E{
			Key:   field,
			Value: 1,
		}
		projection = append(projection, el)
	}

	collection := r.DB.GetConnection().Database(r.dbName).Collection("system.profile")
	ctx := context.TODO()
	cur, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		current := cur.Current

		item := &LogItem{
			Src: r,
		}

		item.RawData = make([]byte, len(current))
		copy(item.RawData, current)
		r.out <- item
	}

	if err := cur.Err(); err != nil {
		return fmt.Errorf("failed to fetch documents from %s", r.DB.GetHost())
	}

	return fmt.Errorf("logLevel is not set to level 1 or 2 on %s", r.DB.GetHost())
}

// Set handler for receiving a logs
func (r *LogsReporter) SetOnDataHandler(handler func(data []byte)) {
	r.onDataHandler = handler
}

// Returns current profiling level
func (r *LogsReporter) GetProfileLevel() (int, error) {
	cmd := bson.D{{
		"profile",
		-1,
	}}

	result := r.DB.RunCommand(context.TODO(), cmd)
	err := result.Err()
	if err != nil {
		return -1, err
	}

	profile := &Profile{}
	err = result.Decode(profile)
	return profile.Was, err
}

/*
	Sets profile level

	2 - The profiler collects data for all operations.
	1 - The profiler collects data for operations that take longer than the value of slowms.
	0 - The profiler is off and does not collect any data. This is the default profiler level.

	Docs: https://docs.mongodb.com/manual/reference/method/db.setProfilingLevel/
 */
func (r *LogsReporter) SetProfileLevel(level int) (int, error) {
	cmd := bson.D{{
		"profile",
		level,
	}}

	result := r.DB.RunCommand(context.TODO(), cmd)
	err := result.Err()
	if err != nil {
		return 0, err
	}

	profile := &Profile{}
	err = result.Decode(profile)
	return profile.Was, err
}