package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup() {
	// get root path
	path, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}

	// Set root path
	racoon.RootPath = path

	// check for .env file
	if !fileExists(path + "/.env") {
		noDotEnv = true
	} else {
		noDotEnv = false
		err = godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		// set db type
		racoon.DB.DatabaseType = os.Getenv("DATABASE_TYPE")
	}
}

func showHelp() {
	helpText := `Available commands:
| command           | args          | description                                                                   |
| :-----------------| :-------------| :-----------------------------------------------------------------------------|
| help              |               | show help text                                                                |
| version           |               | show version                                                                  |
| make new          | appname       | creates a new skeleton app                                                    |
| make migration    | name          | creates 2 new up and down migrations                                          |
| migrate           |               | runs all non run up and down migrations                                       |
| migrate           | down          | reverse the most recent migration                                             |
| migrate           | down force    | force down migration                                                          |
| migrate           | reset         | runs all down migrations in reverse order then all up migrations              |
| make auth         |               | creates and runs migrations for auth tables and creates models and middleware |
| make session      |               | creates a table in the db to use as session store                             |
| make handler      | name          | creates a stub handler in the handlers dir                                    |
| make model        | name          | creates a new model in the data dir                                           |
| make model        | --migrate     | use the migrate flag to also create the db table                              |

`

	if noDotEnv {
		color.Yellow("Warning: no .env file found in current directory, DB functions won't work\n\n")
	}

	color.HiWhite(helpText)
}

func getDSN() (dsn string) {
	dbType := racoon.DB.DatabaseType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	switch dbType {
	case "postgres":
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_POST"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_POST"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
	case "sqlite":
		dsn = fmt.Sprintf("sqlite3://%s/tmp/db-data/%s.db", racoon.RootPath, os.Getenv("DATABASE_NAME"))
	default:
		//
	}

	return dsn
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	// check for an error before doing anything else
	if err != nil {
		return err
	}

	// check if file is directory
	if fi.IsDir() {
		return nil
	}

	// only check go files
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	if matched {
		// read file
		read, err := os.ReadFile(path)
		if err != nil {
			exitGracefully(err)
		}

		// replace placeholder app name and write new file
		updated := strings.ReplaceAll(string(read), "racoonapp", appURL)

		err = os.WriteFile(path, []byte(updated), 0644)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}

func updateSource() {
	// walk entire project folder, including subfolders
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		exitGracefully(err)
	}
}
