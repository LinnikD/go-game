package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Config is our config structure type
type Config struct {
	Token      string `yaml:"token"`
	Mongo      string `yaml:"mongo"`
	Yandex     string `yaml:"yandex"`
}

// LoadConfig retrive information from yaml file
func LoadConfig(path string) (Config, error) {
	var conf Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("Couldn't open file %v: %v", path, err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return conf, fmt.Errorf("Couldn't unmarshal file %v: %v", path, err)
	}
	return conf, nil
}
