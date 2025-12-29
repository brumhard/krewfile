package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const commentSymbol = "#"

type PluginMap map[string]struct{}
type IndexMap map[string]string

type InvalidKrewfileLineError struct {
	line string
}

func (e InvalidKrewfileLineError) Error() string {
	return fmt.Sprintf("parsing krew file, invalid line: %q", e.line)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	var defaultKrewfileLocation string = os.Getenv("KREWFILE")
	if defaultKrewfileLocation == "" {
		defaultKrewfileLocation = filepath.Join(home, ".krewfile")
	}

	var command string
	var krewfileLocation string
	var upgrade, dryRun bool
	flag.StringVar(&command, "command", "krew", "command to be used for krew")
	flag.StringVar(&krewfileLocation, "file", defaultKrewfileLocation, "location of the krewfile")
	flag.BoolVar(&upgrade, "upgrade", false, "runs krew upgrade after syncing plugins")
	flag.BoolVar(&dryRun, "dry-run", false, "shows only output, doesn't modify anything")
	flag.Parse()

	fmt.Println("reading krewfile")
	krewfile, err := os.ReadFile(krewfileLocation)
	if err != nil {
		return err
	}
	wantedPlugins, wantedIndexes, err := readKrewfile(krewfile)
	if err != nil {
		return err
	}

	fmt.Println("getting installed indexes")
	krewIndexList, err := runKrewCommand(false, command, "index", "list")
	if err != nil {
		return err
	}
	installedIndexes := readIndexesFromKrew(krewIndexList)

	// find indexes that are installed but not wanted anymore
	// find indexes that are not installed but wanted
	// -> remove all entries from both maps that are in both maps
	for k := range installedIndexes {
		if _, ok := wantedIndexes[k]; ok {
			delete(installedIndexes, k)
			delete(wantedIndexes, k)
		}
	}

	// now installedIndexes only holds indexes that are not wanted anymore -> to be deleted
	// wantedIndexes are the ones that need to be installed
	for index := range installedIndexes {
		// default index is special case, it's always wanted
		if index == "default" {
			continue
		}
		fmt.Printf("removing index %q\n", index)
		if _, err := runKrewCommand(dryRun, command, "index", "remove", index); err != nil {
			return err
		}
	}

	for index, url := range wantedIndexes {
		fmt.Printf("adding index %q\n", index)
		if _, err := runKrewCommand(dryRun, command, "index", "add", index, url); err != nil {
			return err
		}
	}

	fmt.Println("updating krew")
	if _, err := runKrewCommand(dryRun, command, "update"); err != nil {
		return err
	}

	fmt.Println("getting installed plugins")
	krewPluginList, err := runKrewCommand(false, command, "list")
	if err != nil {
		return err
	}
	installedPlugins := readPluginsFromKrew(krewPluginList)

	// find plugins that are installed but not wanted anymore
	// find plugins that are not installed but wanted
	// -> remove all entries from both maps that are in both maps
	for k := range installedPlugins {
		if _, ok := wantedPlugins[k]; ok {
			delete(installedPlugins, k)
			delete(wantedPlugins, k)
		}
	}

	// now installedPlugins only holds plugins that are not wanted anymore -> to be deleted
	// wantedPlugins are the ones that need to be installed
	for plugin := range installedPlugins {
		fmt.Printf("removing plugin %q\n", plugin)
		pluginItems := strings.Split(plugin, "/")
		if _, err := runKrewCommand(dryRun, command, "uninstall", pluginItems[len(pluginItems)-1]); err != nil {
			return err
		}
	}

	for plugin := range wantedPlugins {
		fmt.Printf("installing plugin %q\n", plugin)
		if _, err := runKrewCommand(dryRun, command, "install", plugin); err != nil {
			return err
		}
	}

	if upgrade {
		fmt.Println("upgrading plugins")
		if _, err := runKrewCommand(dryRun, command, "upgrade"); err != nil {
			return err
		}
	}

	return nil
}

func readKrewfile(input []byte) (PluginMap, IndexMap, error) {
	pluginMap := PluginMap{}
	indexMap := IndexMap{}

	for _, line := range strings.Split(string(input), "\n") {
		lineWithComments := strings.SplitN(line, commentSymbol, 2)
		if parsedLine := strings.TrimSpace(lineWithComments[0]); parsedLine != "" {
			items := strings.Fields(parsedLine)
			if len(items) == 1 {
				pluginMap[items[0]] = struct{}{}
				continue
			}

			if len(items) == 3 && items[0] == "index" {
				indexMap[items[1]] = items[2]
				continue
			}

			return nil, nil, InvalidKrewfileLineError{line: line}
		}

	}

	return pluginMap, indexMap, nil
}

func readPluginsFromKrew(input []byte) PluginMap {
	pluginMap := PluginMap{}

	for _, line := range strings.Split(string(input), "\n") {
		items := strings.Fields(line)
		if len(items) < 1 || items[0] == "PLUGIN" {
			continue
		}

		pluginMap[items[0]] = struct{}{}
	}

	return pluginMap
}

func readIndexesFromKrew(input []byte) IndexMap {
	indexMap := IndexMap{}

	for _, line := range strings.Split(string(input), "\n") {
		items := strings.Fields(line)

		if len(items) < 2 || items[0] == "INDEX" {
			continue
		}

		indexMap[items[0]] = items[1]
	}

	return indexMap
}

func runKrewCommand(dryRun bool, krewCommand string, args ...string) ([]byte, error) {
	fullCommand := append(strings.Split(krewCommand, " "), args...)
	if dryRun {
		fmt.Printf("will run: %q\n", strings.Join(fullCommand, " "))
		return nil, nil
	}

	cmd := exec.Command(fullCommand[0], fullCommand[1:]...)
	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("krew error:\n%s", stderr.String())
	}
	return out, nil
}
