package faraway_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Port       string `json:"port"`
	Difficulty int    `json:"difficulty"`
}

func LoadConfig() Config {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error unmarshalling config:", err)
	}

	return config
}
