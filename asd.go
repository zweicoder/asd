package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zweicoder/asd/config"
	"github.com/zweicoder/asd/utils"
	// "log"
	"os"
	"os/exec"
	"path/filepath"
)

func updateCache(c *cli.Context) error {
	// Download listed modules, analyze if they have deps, recursively install deps (and check in a set to make sure no circular)
	fmt.Printf("Downloading modules to %s...\n", config.CachePath)
	moduleZip := "/tmp/asd-modules.zip"
	err := utils.DownloadFile(moduleZip, config.ModuleURL)
	if err != nil {
		return err
	}

	err = utils.Unzip(moduleZip, config.CachePath)
	if err != nil {
		return err
	}

	err = os.Remove(moduleZip)
	if err != nil {
		return err
	}
	return nil
}

func generateCommands(c *cli.Context) (string, error) {
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, "install")
		return "", cli.NewExitError("Error: Nothing to install!", 1)
	}

	// Create cache if needed
	if _, err := os.Stat(config.CachePath); os.IsNotExist(err) {
		// Download to .asd cache, maybe even version the instruction sets
		err := os.MkdirAll(config.CachePath, 0755)
		if err != nil {
			return "", cli.NewExitError(err, 1)
		}
		err = updateCache(c)
		if err != nil {
			return "", err
		}
	}

	// Parse all requested deps into a list of instructions
	args := c.Args()
	arr := getDeps(args)

	var buffer bytes.Buffer
	if c.String("remote") != "" {
		remote := c.String("remote")
		buffer.WriteString(fmt.Sprintf("ssh %v 'bash -s' << EOF \n", remote))
	} else {
		buffer.WriteString("bash -s << EOF \n")
	}
	var keys []string
	for _, e := range arr {
		keys = append(keys, e.key)
	}
	fmt.Printf("Installing modules: %v\n", keys)

	for _, e := range arr {
		for _, command := range e.commands {
			buffer.WriteString(command)
			buffer.WriteString("\n")
		}
	}
	buffer.WriteString("EOF\n")
	return buffer.String(), nil
}

// CliUpdateCache updates the asd-modules cache
func CliUpdateCache(c *cli.Context) error {
	e := updateCache(c)
	if e != nil {
		return cli.NewExitError(e, 1)
	}
	return nil
}

// CliGen generates the install script. More for debugging as of now
func CliGen(c *cli.Context) error {
	commands, e := generateCommands(c)
	if e != nil {
		return e
	}

	var f *os.File
	if pathFlag := c.String("path"); pathFlag != "" {
		path, e := filepath.Abs(pathFlag)
		if e != nil {
			return e
		}
		f, e = os.Create(path)
	} else {
		f, e = os.Create("install.sh")
	}
	defer f.Close()
	f.WriteString(commands)
	return nil
}

// Version 0.1.0 of this assumes the git repository is cloned (with module files inside).
// This function just tries to do some simple dep resolution and install the files in the repo
func CliInstall(c *cli.Context) error {
	commands, e := generateCommands(c)
	if e != nil {
		return e
	}

	// log.Printf("commands to run: %v\n", commands)
	if e != nil {
		return e
	}

	var f *os.File
	if pathFlag := c.String("path"); pathFlag != "" {
		path, e := filepath.Abs(pathFlag)
		if e != nil {
			return e
		}
		f, e = os.Create(path)
	} else {
		f, e = os.Create("/tmp/install.sh")
	}
	defer f.Close()
	f.WriteString(commands)
	cmd := exec.Command("bash", "/tmp/install.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e = cmd.Run()
	if e != nil {
		return e
	}
	return nil
}
