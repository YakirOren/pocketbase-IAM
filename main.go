package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "pocketbase-iam/migrations"

	"pocketbase-iam/iam"
)

func main() {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	iam.RegisterRoutes(app)
	iam.RegisterHooks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
