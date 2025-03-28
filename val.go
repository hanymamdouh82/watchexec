package main

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	cooldownTime   = 500 * time.Millisecond // Short time to check for quick deletes
	fileTimestamps = make(map[string]time.Time)
	fileLock       sync.Mutex
)

// decides which dirs to be ignored during early initiation of watched dirs
func shouldIgnoreDir(path string, rootDir string) bool {
	config.mu.Lock()
	defer config.mu.Unlock()

	dir := config.Dirs[rootDir]
	for _, pattern := range dir.Exclude {
		if strings.Contains(path, pattern) { // Simple substring match
			return true
		}
	}
	return false
}

// decides if the event for file should be ignored from excuting the bin or not
func shouldIgnoreFile(path string) bool {
	config.mu.Lock()

	var rootDir string
	for k := range config.Dirs {
		if strings.Contains(path, k) { // Simple substring match
			rootDir = k
			break
		}
	}
	if rootDir == "" { // Prevent nil access
		config.mu.Unlock()
		return false
	}

	excludeList := config.Dirs[rootDir].Exclude
	config.mu.Unlock()

	for _, pattern := range excludeList {
		match, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && match {
			return true
		}
	}

	return false
}

// Track files before triggering the function
func shouldTrigger(event fsnotify.Event) bool {
	fileLock.Lock()
	defer fileLock.Unlock()

	// Ignore excluded files
	if shouldIgnoreFile(event.Name) {
		return false
	}

	switch event.Op {
	case fsnotify.Create:
		// Store the creation timestamp
		fileTimestamps[event.Name] = time.Now()
		go watched.CheckForNewDirs(event.Name) // Check if it's a new directory to watch
		return false                           // Don't trigger immediately

	case fsnotify.Remove:
		// Check if the file existed briefly
		if _, exists := fileTimestamps[event.Name]; exists {
			delete(fileTimestamps, event.Name) // Clean up
			return false                       // Ignore quick create-delete events
		}
		return true // Normal delete event

	case fsnotify.Write, fsnotify.Chmod:
		// Ignore temporary changes if the file was recently created and removed
		if _, exists := fileTimestamps[event.Name]; exists {
			return false
		}
		return true // Allow normal writes

	default:
		return true // Trigger for other cases
	}
}
