package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type JobConfig struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name  string `yaml:"name"`
	Shell string `yaml:"shell"`
	Run   string `yaml:"run"`
}

type Event struct {
	Tables []string `yaml:"tables"`
}

type Config struct {
	Name string            `yaml:"name"`
	On   map[string]*Event `yaml:"on"`
	Jobs []JobConfig       `yaml:"jobs"`
}

func ReadConfigFile(configFile string) Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	return config
}
