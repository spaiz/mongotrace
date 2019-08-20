package app

import (
	"context"
	"github.com/spaiz/mongotrace/config"
	"github.com/spaiz/mongotrace/db"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sync"
)

func NewSetCommand(conf *config.Configuration, options *Options, dbs []*db.MongoDB) *SetLevelCommand {
	return &SetLevelCommand{
		conf:    conf,
		options: options,
		dbs:     dbs,
	}
}

type SetLevelCommand struct {
	options       *Options
	previousLevel int
	conf          *config.Configuration
	dbs           []*db.MongoDB
}

func (r *SetLevelCommand) Execute() error {
	hostsNum := len(r.dbs)
	var wg sync.WaitGroup
	wg.Add(hostsNum)

	for _, dbconn := range r.dbs {
		go func(dbd *db.MongoDB) {
			defer wg.Done()
			cmd := bson.D{{
				"profile",
				r.options.Level,
			}}

			result := dbd.RunCommand(context.TODO(), cmd)
			err := result.Err()
			if err != nil {
				log.Fatal(err)
			}

			profile := &Profile{}
			err = result.Decode(profile)
			r.previousLevel = profile.Was

			if r.options.Level == 2 {
				log.Printf("Host: %s, Status: enabled\n", dbd.GetHost())
			} else if r.options.Level == 0 {
				log.Printf("Host: %s, Status: disabled\n", dbd.GetHost())
			}
		}(dbconn)
	}

	wg.Wait()

	return nil
}

func (r *SetLevelCommand) setLevel() error {
	return nil
}
