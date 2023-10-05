package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

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

	var command string
	var krewfileLocation string
	var upgrade bool
	flag.StringVar(&command, "command", "krew", "command to be used for krew")
	flag.StringVar(&krewfileLocation, "file", filepath.Join(home, ".krewfile"), "location of the krewfile")
	flag.BoolVar(&upgrade, "upgrade", false, "runs krew upgrade after syncing plugins")
	flag.Parse()

	fmt.Println("getting installed plugins")
	krewList, err := krewCommand(command, "list").Output()
	if err != nil {
		return err
	}

	fmt.Println("reading krewfile")
	krewfile, err := os.ReadFile(krewfileLocation)
	if err != nil {
		return err
	}

	installedPlugins := readBytesToPluginMap(krewList)
	wantedPlugins := readBytesToPluginMap(krewfile)

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
		fmt.Printf("removing %s\n", plugin)
		if err := krewCommand(command, "uninstall", plugin).Run(); err != nil {
			return err
		}
	}

	fmt.Println("updating krew")
	if err := krewCommand(command, "update").Run(); err != nil {
		return err
	}

	for plugin := range wantedPlugins {
		fmt.Printf("installing %s\n", plugin)
		if err := krewCommand(command, "install", plugin).Run(); err != nil {
			return err
		}
	}

	if upgrade {
		fmt.Println("upgrading plugins")
		if err := krewCommand(command, "upgrade").Run(); err != nil {
			return err
		}
	}

	return nil
}

func readBytesToPluginMap(input []byte) map[string]struct{} {
	output := map[string]struct{}{}
	for _, line := range strings.Split(string(input), "\n") {
		if line != "" {
			output[strings.TrimSpace(line)] = struct{}{}
		}
	}

	return output
}

func krewCommand(krewCommand string, args ...string) *exec.Cmd {
	fullCommand := append(strings.Split(krewCommand, " "), args...)
	return exec.Command(fullCommand[0], fullCommand[1:]...)
}
