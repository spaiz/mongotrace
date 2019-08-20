package db

import (
	"context"
	"fmt"
	"github.com/spaiz/mongotrace/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoDBConfig struct {
	ConnectionString string
}

func NewMongoDB(conf *config.HostConfig, timeout time.Duration) *MongoDB {
	return &MongoDB{
		conf: conf,
		timeout: timeout,
	}
}

type MongoDB struct {
	conf    *config.HostConfig
	client  *mongo.Client
	opts    *options.ClientOptions
	timeout time.Duration
}

func (r *MongoDB) Connect() error {
	var err error
	r.opts = options.Client().ApplyURI(r.conf.ConnectionString)
	r.client, err = mongo.NewClient(r.opts)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), r.timeout)
	err = r.client.Connect(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoDB) GetConnection() *mongo.Client {
	return r.client
}

func (r *MongoDB) Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), r.timeout)
	return r.client.Ping(ctx, readpref.Primary())
}
func (r *MongoDB) Name() string {
	return r.conf.DBName
}

func (r *MongoDB) GetHost() string {
	return fmt.Sprintf("%v", r.opts.Hosts)
}

func (r *MongoDB) RunCommand(ctx context.Context, cmd bson.D) *mongo.SingleResult {
	return r.GetConnection().Database(r.Name()).RunCommand(ctx, cmd)
}
