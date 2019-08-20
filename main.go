package main // import "github.com/spaiz/mongotrace"

import (
	"github.com/spaiz/mongotrace/app"
	"log"

)

var (
	// VERSION will be set during build
	VERSION = "0.0.0"
)

func main() {
	app := app.NewApp(VERSION)
	if err := app.Run(); err != nil {
		log.Fatalf("Exited with error: %s\n", err.Error())
	}
}
