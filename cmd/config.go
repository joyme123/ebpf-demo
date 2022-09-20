package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/fsnotify.v1"
	"gopkg.in/yaml.v2"
)

// Config holds the parsed configuration file
type Config struct {
	Links map[string]string `yaml:"links"`
}

func newFromFile(file string) (*Config, error) {
	var cfg Config
	log.Infof("Parsing config file: %s", file)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	log.Infof("App configuration: %+v", cfg)
	return &cfg, nil
}

func configWatcher(file string, out chan map[string]string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicf("Failed to initialize fsnotify: %s", err)
	}

	defer watcher.Close()

	if err := watcher.Add(file); err != nil {
		log.Panicf("Error watching the configuration file: %s", err)
	}

	for {
		select {
		case <-watcher.Events:
			cfg, err := newFromFile(file)
			if err != nil {
				log.Infof("Error parsing the configuration file: %s", err)
			} else {
				out <- cfg.Links
			}

		case err := <-watcher.Errors:
			log.Infof("Received watcher.Error: %s", err)
		}
	}
}
