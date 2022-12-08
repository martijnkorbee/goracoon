package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

func doMakeSession() error {
	dbType := racoon.DB.DatabaseType
	if dbType == "postgresql" {
		dbType = "postgres"
	}
	fileName := fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())
	upFile := racoon.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := racoon.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/sessions."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	switch dbType {
	case "postgres":
		err = os.WriteFile(downFile, []byte("drop index if exists sessions_expiry_idx cascade; drop table if exists sessions cascade;"), 0644)
		if err != nil {
			exitGracefully(err)
		}
	case "sqlite":
		err = os.WriteFile(downFile, []byte("drop index if exists sessions_expiry_idx; drop table if exists sessions;"), 0644)
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

	color.White("INFO: sessions table created in db: %s", dbType)
	color.Green("Status: OK")

	return nil
}
