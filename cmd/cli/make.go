package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3 string) error {
	switch arg2 {
	case "key":
		key := racoon.RandomStringGenerator(32)
		white := color.New(color.FgHiWhite).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		fmt.Printf("%s: %s\n", blue("32 character encryption key"), white(key))

	case "migration":
		// get database type
		dbType := racoon.DB.DatabaseType

		if arg3 == "" {
			exitGracefully(errors.New("you must give the migration a name"))
		}

		// create migration file name
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := racoon.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := racoon.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}

	case "auth":
		err := doMakeAuth()
		if err != nil {
			exitGracefully(err)
		}

	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the new handler a name"))
		}

		fileName := racoon.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}

		color.White("INFO: handler created: " + fileName)
		color.Green("Status: OK")

	case "model":
		// check if name arg is passed
		if arg3 == "" {
			exitGracefully(errors.New("you must give the model a name"))
		}

		// extract make model flags
		f := flag.NewFlagSet("make model <name>", flag.ContinueOnError)
		doMigrate := f.Bool("migrate", false, "create model migrations")
		f.Parse(os.Args[4:])

		err := doMakeModel(arg3, *doMigrate)
		if err != nil {
			exitGracefully(err)
		}

	case "session":
		err := doMakeSession()
		if err != nil {
			exitGracefully(err)
		}

	default:
		// default
	}

	return nil
}
