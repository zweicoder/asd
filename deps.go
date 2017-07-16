package main

import (
	"fmt"
	"github.com/zweicoder/asd/config"
	"gopkg.in/fatih/set.v0"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type dep struct {
	key  string
	deps []dep
}

type conf struct {
	Module struct {
		Commands     []string
		Script       string
		Dependencies []string
	}
}

// ModuleNode stores information on what command(s) to run, although right now it's either
// `bash another-script.sh` or `apt-get something`
type ModuleNode struct {
	key string
	// wrapped command from the struct
	commands []string
	deps     []string
}

func getInfo(key string) (info ModuleNode) {
	// read from yaml
	c := conf{}
	dir := path.Join(config.ModulePath, key)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		command := fmt.Sprintf("sudo apt-get install %s", key)
		return ModuleNode{key: key, commands: []string{command}}
	}

	filepath := path.Join(dir, "module.yml")
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// Attempt to automagically find .sh file if no module.yml file
		filepath = path.Join(dir, fmt.Sprintf("%s.sh", key))
		if _, err = os.Stat(filepath); os.IsNotExist(err) {
			// Directory exists but .yml and .sh not defined. Just panic here
			panic(err)
		}

		command := fmt.Sprintf("`cat %s`", filepath)
		return ModuleNode{key: key, commands: []string{command}}
	}

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(file, &c); err != nil {
		// Probably bad Yaml config
		panic(err)
	}

	if c.Module.Script != "" && len(c.Module.Commands) > 0 {
		log.Fatalf("Both 'script' and 'commands' are defined in %s\n", filepath)
	} else if c.Module.Script == "" && len(c.Module.Commands) == 0 {
		log.Fatalf("Bad config for %s\n", filepath)
	}

	if c.Module.Script != "" {
		command := fmt.Sprintf("`cat %s`", path.Join(dir, c.Module.Script))
		return ModuleNode{
			key:      key,
			commands: []string{command},
			deps:     c.Module.Dependencies,
		}
	}
	return ModuleNode{
		key:      key,
		commands: c.Module.Commands,
		deps:     c.Module.Dependencies,
	}
}

func resolveDep(item string, ret *[]ModuleNode, explored *set.Set) {
	if explored.Has(item) {
		return
	}

	info := getInfo(item)
	explored.Add(item)
	// DFS by recursively following other edges before adding to return array
	for _, dep := range info.deps {
		resolveDep(dep, ret, explored)
	}
	*ret = append(*ret, info)
}

// Get an array of flattened deps recursively
func getDeps(items []string) (ret []ModuleNode) {
	explored := set.New()
	ret = []ModuleNode{}
	var tail string
	// Explore each vertex
	for len(items) > 0 {
		tail, items = items[len(items)-1], items[:len(items)-1]
		resolveDep(tail, &ret, explored)
	}

	return
}
