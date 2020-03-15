package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TslConfig struct {
	DefaultFrom string `json:"from"`
	DefaultTo   string `json:"to"`
	ExampleNeed bool   `json:"expneed"`
}

var (
	tslBaseLoc   = filepath.Join(os.Getenv("HOME"), ".tsl")
	tslDbLoc     = filepath.Join(tslBaseLoc, "db")
	tslConfigLoc = filepath.Join(tslBaseLoc, "tsl_config.json")
)

func writeTslConfig(cfg *TslConfig) error {
	if cfg == nil {
		return errors.New("TslConfig is nil")
	}

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(cfg); err != nil {
		return err
	}

	return ioutil.WriteFile(tslConfigLoc, buf.Bytes(), 0644)
}

func loadTslConfig() (*TslConfig, error) {
	f, err := os.OpenFile(tslConfigLoc, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var cfg TslConfig

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewTslConfig() (*TslConfig, error) {
	cfg, err := loadTslConfig()
	if err == nil {
		return cfg, err
	}

	if _, err := os.Stat(tslBaseLoc); err != nil {
		if err := os.Mkdir(tslBaseLoc, 0755); err != nil {
			return nil, err
		}
	}

	//TODO: initialize DB

	cfg = &TslConfig{DefaultFrom: "en", DefaultTo: "zh-Hans", ExampleNeed: false}
	if err := writeTslConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
