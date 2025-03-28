package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"flag"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

var (
	verbose bool
	config  Config
	watcher *fsnotify.Watcher
	watched *WatchManager
)

func init() {
	var confFile string
	flag.StringVar(&confFile, "c", "conf.yml", "configuation file name (.yml file)")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	data, err := os.ReadFile(confFile)
	if err != nil {
		log.Fatal("failed to load conf.yaml, exiting...!")
	}

	config.mu.Lock()
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatal("Failed to parse conf.yaml, exiting...!")
	}
	config.mu.Unlock()

	validateConfig()

	config.mu.Lock()
	if verbose {
		LogJSON(config.Dirs)
	}
	config.mu.Unlock()
}

func validateConfig() {
	config.mu.Lock()
	defer config.mu.Unlock()

	for dir, dconf := range config.Dirs {
		if len(dir) == 0 {
			log.Fatal("empty dir name, exiting...!")
		}

		// validate if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			logMsg(fmt.Sprintf(dir, "does not exist, ignoring dir"))
			delete(config.Dirs, dir)
		}

		// ensure provided bin to excute
		// ignore if bin is not found
		if len(dconf.Bin) == 0 {
			delete(config.Dirs, dir)
		}
	}

	// Exit if directories are empty after running the validate
	if len(config.Dirs) == 0 {
		log.Fatal("Empty directories to watch...!")
	}
}

func main() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("error creating watcher: ", err)
	}
	defer watcher.Close()

	// Watch all initial directories recursively
	watched = NewWatchManager()
	config.mu.Lock()
	// copy dirs to new variable and Unlock
	// this is because WatchDirRecursive also locks the config variable
	dirs := config.Dirs
	config.mu.Unlock()

	for dir := range dirs {
		if err := watched.WatchDirRecursive(dir); err != nil {
			logMsg(fmt.Sprintf("Error watching: %s, %v", dir, err))
		}
	}

	// Process changes
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			dir := filepath.Dir(event.Name)

			// Check if this event should be ignored
			if !shouldTrigger(event) {
				continue
			}

			// handle the required bin
			triggerBinDebounced(dir)

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logMsg(fmt.Sprintf("Watcher error: %v", err))
		}
	}
}

func LogJSON(val any) {
	j, _ := json.MarshalIndent(val, "", "    ")
	fmt.Printf("%s\n", j)
}

func logMsg(msg string) {
	if !verbose {
		return
	}

	fmt.Println(msg)
}
