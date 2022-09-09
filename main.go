package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_ "go.uber.org/automaxprocs"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("getting installed plugins")
	krewList, err := exec.Command("krew", "list").Output()
	if err != nil {
		return err
	}

	fmt.Println("reading krewfile")
	krewfile, err := os.ReadFile("./krewfile")
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
		if err := exec.Command("krew", "uninstall", plugin).Run(); err != nil {
			return err
		}
	}

	fmt.Println("updating krew")
	if err := exec.Command("krew", "update").Run(); err != nil {
		return err
	}

	for plugin := range wantedPlugins {
		fmt.Printf("installing %s\n", plugin)
		if err := exec.Command("krew", "install", plugin).Run(); err != nil {
			return err
		}
	}

	fmt.Println("upgrading plugins")
	if err := exec.Command("krew", "upgrade").Run(); err != nil {
		return err
	}

	return nil
}

func readBytesToPluginMap(input []byte) map[string]struct{} {
	output := map[string]struct{}{}
	for _, line := range strings.Split(string(input), "\n") {
		if line != "" {
			output[line] = struct{}{}
		}
	}

	return output
}
