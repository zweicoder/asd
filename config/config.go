package config

import (
	"os"
	"path"
)

// HomeDir is just the homeDir in the user's system
var HomeDir string = os.Getenv("HOME")

// CachePath is where the modules will be downloaded
var CachePath string = path.Join(HomeDir, ".asd")
var ModulePath string = path.Join(CachePath, "asd-modules-master")

// ModuleURL is the repo containing the collection of scripts
const ModuleURL string = "https://github.com/zweicoder/asd-modules/archive/master.zip"
