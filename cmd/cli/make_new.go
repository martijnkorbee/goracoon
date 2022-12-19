package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

var appURL string

func doNew(appName string) error {
	appName = strings.ToLower(appName)
	appURL = appName

	// sanitize the application name (convert url to single word)
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[len(exploded)-1]
	}

	// git clone the skeleton application
	color.Yellow("\tCloning git repository...")

	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "git@github.com:martijnkorbee/racoon-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a ready to go .env file
	color.Yellow("\tCreating .env file...")

	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		exitGracefully(err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", racoon.RandomStringGenerator(32))

	err = os.WriteFile(fmt.Sprintf("./%s/.env", appName), []byte(env), 0644)
	if err != nil {
		exitGracefully(err)
	}

	color.HiWhite("Finished creating .env file")

	// update the go.mod file
	color.Yellow("\tCreating go.mod file...")

	_ = os.Remove("./" + appName + "/go.mod")

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = os.WriteFile(fmt.Sprintf("./%s/go.mod", appName), []byte(mod), 0644)
	if err != nil {
		exitGracefully(err)
	}

	color.HiWhite("Finished creating go.mod file")

	// change to new project dir
	err = os.Chdir("./" + appName)
	if err != nil {
		exitGracefully(fmt.Errorf("couldn't change to new app dir: %s", err))
	}

	// update the existing .go files with correct name and imports
	color.Yellow("\tUpdating source files...")
	updateSource()
	color.HiWhite("Updated source files")

	// run go mod tidy in the project directory
	color.Yellow("\tRunning go mod tidy...")

	cmd := exec.Command("go", "mod", "tidy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		exitGracefully(err)
	}

	color.HiWhite("Finished running go mod tidy")

	// TODO: create makefile support for windows
	// update makefile
	color.Yellow("\tUpdating makefile")

	data, err = os.ReadFile("Makefile.unix")
	if err != nil {
		exitGracefully(err)
	}
	makefile := strings.ReplaceAll(string(data), "${APP_NAME}", appName)
	err = os.WriteFile("Makefile", []byte(makefile), 0644)
	if err != nil {
		exitGracefully(err)
	}

	err = os.Remove("Makefile.unix")
	if err != nil {
		color.Yellow("WARNING: could not remove Makefile.unix: ", err)
	}

	color.HiWhite("Finished updating makefile")
	color.Green("Done creating new %s, status: OK", appName)

	color.Yellow("\n\tBuilding: %s", appName)
	cmd = exec.Command("make", "build")
	output, err = cmd.CombinedOutput()
	if err != nil {
		color.Red(fmt.Sprint(err) + ": " + string(output))
		exitGracefully(err)
	}
	color.HiWhite("%s", string(output))

	color.Green("Start your app from dir %s and run: /tmp/%s", appName, appName)

	return nil
}
