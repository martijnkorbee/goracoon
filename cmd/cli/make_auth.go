package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

func doMakeAuth() error {
	// migrations
	dbType := racoon.DB.DatabaseType
	if dbType == "postgresql" {
		dbType = "postgres"
	}

	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := racoon.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := racoon.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	switch dbType {
	case "postgres":
		err = os.WriteFile(downFile, []byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens cascade;"), 0644)
		if err != nil {
			exitGracefully(err)
		}
	case "sqlite":
		err = os.WriteFile(downFile, []byte("drop table if exists users; drop table if exists tokens; drop table if exists remember_tokens;"), 0644)
		if err != nil {
			exitGracefully(err)
		}
	default:
		//
	}

	// run migrations
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	// copy user and token models
	err = copyFileFromTemplate("templates/data/user.go.txt", racoon.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/token.go.txt", racoon.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy auth middleware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt", racoon.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", racoon.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("\nWARNING: must add user and token models in data/models.go and appropriate middleware to your routes")
	color.Green("\nStatus: OK")

	return nil
}
