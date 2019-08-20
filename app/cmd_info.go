package app

import (
	"context"
	"github.com/spaiz/mongotrace/config"
	"github.com/spaiz/mongotrace/db"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sync"
)

// Returns New InfoCommand instance
func NewInfoCommand(conf *config.Configuration, dbs []*db.MongoDB) *InfoCommand {
	return &InfoCommand{
		conf: conf,
		dbs:  dbs,
	}
}

// Defines InfoCommand struct
type InfoCommand struct {
	previousLevel int
	conf          *config.Configuration
	dbs          []*db.MongoDB
}

// Executes the command
func (r *InfoCommand) Execute() error {
	hostsNum := len(r.dbs)

	var wg sync.WaitGroup
	wg.Add(hostsNum)

	for _, dbconn := range r.dbs {
		go func(dbd *db.MongoDB) {
			defer wg.Done()

			cmd := bson.D{{
				"profile",
				-1,
			}}

			result := dbd.RunCommand(context.TODO(), cmd)
			err := result.Err()
			if err != nil {
				log.Fatal(err)
			}

			profile := &Profile{}
			err = result.Decode(profile)
			r.previousLevel = profile.Was
			profile.DbName = dbd.Name()
			profile.Host = dbd.GetHost()

			r.printInfo(profile)

		}(dbconn)
	}

	wg.Wait()

	return nil
}

func (r *InfoCommand) printInfo(profile *Profile) {
	log.Printf("Host: %s, Level: %d\n", profile.Host, profile.Was)
}
