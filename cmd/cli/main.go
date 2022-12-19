package main

import (
	"errors"
	"os"

	"github.com/fatih/color"
	"github.com/martijnkorbee/goracoon"
)

const version = "1.0.0"

var racoon goracoon.Goracoon
var noDotEnv bool

func main() {
	var message string

	setup()

	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}

	switch arg1 {
	case "help":
		showHelp()

	case "version":
		color.Yellow("Application version: " + version)

	case "maintenance":
		if arg2 == "" {
			exitGracefully(errors.New("maintenance requires up or down: maintenance up|down"))
		}

		if noDotEnv {
			exitGracefully(errors.New("no .env file in current directory"))
		}

		switch arg2 {
		case "up":
			rpcMaintenanceMode(true)
		case "down":
			rpcMaintenanceMode(true)
		default:
			//
		}

	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("new requires an application name"))
		}
		err = doNew(arg2)
		if err != nil {
			exitGracefully(err)
		}

	case "migrate":
		if noDotEnv {
			exitGracefully(errors.New("no .env file in current directory"))
		}

		if arg2 == "" {
			arg2 = "up"
		}

		err = doMigrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}

		message = "Migrations complete!"

	case "make":
		if noDotEnv {
			exitGracefully(errors.New("no .env file in current directory"))
		}

		if arg2 == "" {
			exitGracefully(errors.New("make requires a sub command: (migration|model|handler)"))
		}
		err = doMake(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}

	default:
		showHelp()
	}

	exitGracefully(nil, message)
}

func validateInput() (arg1 string, arg2 string, arg3 string, err error) {
	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	} else {
		showHelp()
		return "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, nil
}

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}

	if len(message) > 0 {
		color.Yellow(message)
	}

	os.Exit(0)
}
