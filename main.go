package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	DEFAULT_CONFIG = "/etc/sshauth/config.json"
)

type Config struct {
	Token string
	Owner string
	Team  string
}

func loadConfig(file string) Config {
	f, err := os.Open(file)
	exitIf(err)

	decoder := json.NewDecoder(f)
	config := Config{}
	err = decoder.Decode(&config)
	exitIf(err)

	return config
}

func exitIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var configFile = flag.String("config", DEFAULT_CONFIG, "path to a JSON config file")
	flag.Parse()

	config := loadConfig(*configFile)

	c := NewGithubClient(config.Token, config.Owner)

	users, err := c.GetTeamMembers(config.Team)
	exitIf(err)

	keys := c.GetTeamKeys(users)
	for _, k := range keys {
		fmt.Println(*k.Key)
	}
}
