package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMakeModel(arg3 string, doMigrate bool) error {
	// create a new pluralize instance
	plur := pluralize.NewClient()

	// correct model naming if necassery
	var modelName = arg3
	var tableName = arg3

	if plur.IsPlural(arg3) {
		modelName = plur.Singular(arg3)
		tableName = strings.ToLower(arg3)
	} else {
		tableName = strings.ToLower(plur.Plural(arg3))
	}

	// create model
	createModel(modelName, tableName)

	if doMigrate {
		// create model table
		createModelTable(modelName, tableName)
	}

	// print feedback
	color.Green("Status: OK")

	return nil
}

func createModel(modelName, tableName string) {
	// create target file name (full path)
	fileName := racoon.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
	if fileExists(fileName) {
		exitGracefully(errors.New(fileName + " already exists"))
	}

	// read model go text file
	data, err := templateFS.ReadFile("templates/data/model.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	// update model go text file with model name
	model := string(data)
	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

	err = os.WriteFile(fileName, []byte(model), 0644)
	if err != nil {
		exitGracefully(err)
	}

	// print model feedback
	color.White("INFO: model created: " + fileName)
}

// createModelTable create db table corresponding to the created model
func createModelTable(modelName, tableName string) {
	dbType := racoon.DB.DatabaseType
	if dbType == "postgresql" {
		dbType = "postgres"
	}
	fileName := fmt.Sprintf("%d_create_%s_table", time.Now().UnixMicro(), strings.ToLower(tableName))
	upFile := racoon.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := racoon.RootPath + "/migrations/" + fileName + ".down.sql"

	// read model table go text file
	data, err := templateFS.ReadFile("templates/migrations/model_table." + dbType + ".sql.txt")
	if err != nil {
		exitGracefully(err)
	}

	// update model table go text file with tablename
	table := string(data)
	table = strings.ReplaceAll(table, "$TABLENAME$", tableName)

	err = os.WriteFile(upFile, []byte(table), 0644)
	if err != nil {
		exitGracefully(err)
	}

	switch dbType {
	case "postgres":
		err = os.WriteFile(downFile, []byte(fmt.Sprintf("drop table if exists %s cascade;", tableName)), 0644)
		if err != nil {
			exitGracefully(err)
		}
	case "sqlite":
		err = os.WriteFile(downFile, []byte(fmt.Sprintf("drop table if exists %s;", tableName)), 0644)
		if err != nil {
			exitGracefully(err)
		}
	default:
		//
	}

	// run migration for model
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	// print model table migrate feedback
	color.White("INFO: model table migrated: " + tableName)
}
