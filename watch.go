package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WatchManager manages directory watching with a map and lock
type WatchManager struct {
	mu       sync.Mutex
	watchMap map[string]string
}

func NewWatchManager() *WatchManager {
	return &WatchManager{
		watchMap: make(map[string]string),
	}
}

// Add a directory to watchMap safely
func (wm *WatchManager) AddWatchDir(dir, root string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.watchMap[dir] = root
}

// Check if a directory is already being watched
func (wm *WatchManager) IsWatched(dir string) bool {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	_, exists := wm.watchMap[dir]
	return exists
}

func (wm *WatchManager) GetRoot(dir string) string {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	root, exists := wm.watchMap[dir]
	if exists {
		return root
	}
	return ""
}

// for debugging purposes only
func (wm *WatchManager) Log() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	LogJSON(wm.watchMap)
}

// Recursively watch directories and add to watchMap
func (wm *WatchManager) WatchDirRecursive(rootDir string) error {
	wm.AddWatchDir(rootDir, rootDir) // Ensure rootDir is in watchMap

	err := filepath.Walk(rootDir, func(subDir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// if !info.IsDir() || shouldIgnore(subDir) {
		if !info.IsDir() || shouldIgnoreDir(subDir, rootDir) {
			return nil
		}

		// Add subdir under rootDir
		wm.AddWatchDir(subDir, rootDir)

		// Now add it to the fsnotify watcher (without locking)
		err = watcher.Add(subDir)
		if err != nil {
			logMsg(fmt.Sprintf("Error watching dir: %s | %v", subDir, err))
		} else {
			logMsg(fmt.Sprintf("Watching: %s", subDir))
		}

		return nil
	})

	return err
}

// Detect newly created directories and start watching them
func (wm *WatchManager) CheckForNewDirs(path string) {
	time.Sleep(50 * time.Millisecond) // Give FS some time to settle

	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		if !wm.IsWatched(path) {
			wm.WatchDirRecursive(path)
		}
	}
}
