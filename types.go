package main

import "sync"

type DirConfig struct {
	Bin           string   `yaml:"bin"`
	Args          []string `yaml:"args"`
	Stdout        bool     `yaml:"stdout"`
	Exclude       []string `yaml:"exclude"`
	DebounceDelay int      `yaml:"delay"`
}

type Config struct {
	mu   sync.Mutex
	Dirs map[string]DirConfig `yaml:"dirs"`
}
