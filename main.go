package main

import (
	"log"

	"github.com/pocketbase/pocketbase"

	"github.com/yakiroren/pocketbase-IAM/iam"
)

func main() {
	app := pocketbase.New()

	if err := iam.Setup(app, iam.DefaultOptions()); err != nil {
		log.Fatalf("Failed to setup IAM: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
