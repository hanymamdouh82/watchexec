package main

import (
	"fmt"
	"os/exec"
	"time"
)

var timers = make(map[string]*time.Timer)

// Debounce sync execution to avoid rapid Git commits
func triggerBinDebounced(dir string) {

	root := watched.GetRoot(dir)

	fileLock.Lock()
	if timer, exists := timers[root]; exists {
		timer.Stop()
	}
	fileLock.Unlock()

	config.mu.Lock()
	debounceDelay := time.Second * time.Duration(config.Dirs[root].DebounceDelay)
	config.mu.Unlock()
	// we dont accept immediate debounce and we set to 1 second
	if debounceDelay == 0 {
		debounceDelay = time.Second
	}
	timers[root] = time.AfterFunc(debounceDelay, func() {
		execBin(dir)
	})
}

func execBin(dir string) {
	root := watched.GetRoot(dir)

	config.mu.Lock()
	bin := config.Dirs[root].Bin
	args := config.Dirs[root].Args
	stdout := config.Dirs[root].Stdout
	config.mu.Unlock()

	cmd := exec.Command(bin, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error excuting command: %s | error: %v\n", bin, err)
	}

	if stdout {
		fmt.Printf("%s", string(output))
	}
}
