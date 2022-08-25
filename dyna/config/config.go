package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DBConnection string `json:"db_connection"`
	Port         string `json:"port"`
	RedisAddr    string `json:"redis_addr"`
}

// returns a new config from filename passed
func NewConfig(fileName string) (*Config, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
