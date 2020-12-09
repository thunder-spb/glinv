package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config model
type Config struct {
	Addr   string `json:"addr"`
	DSN    string `json:"dsn"`
	Secret string `json:"secret"`
}

// LoadConfiguration ...
func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer func() {
		cerr := configFile.Close()

		if err == nil {
			err = cerr
		}
	}()

	if err != nil {
		log.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		log.Fatal(err)
	}

	return config
}
