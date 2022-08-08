package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hazelcast/hazelcast-go-client"
)

type Config struct {
	Client         hazelcast.Config
	OperationCount int
	Concurrency    int
}

func LoadConfigFromPath(path string) (*Config, error) {
	if filepath.Ext(path) != ".json" {
		return nil, fmt.Errorf("cannot load config from path: %s: only JSON files are supported", path)
	}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("path not found: %s: %w", path, err)
		}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %s: %w", path, err)
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}
	if c.OperationCount == 0 {
		c.OperationCount = 1000
	}
	if c.Concurrency == 0 {
		c.Concurrency = 1
	}
	return &c, nil
}
