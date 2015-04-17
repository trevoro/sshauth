package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

const defaultConfig = "/etc/sshauth/config.json"

type config struct {
	Token string
	Owner string
	Team  string
}

func loadConfig(file string) config {
	f, err := os.Open(file)
	exitIf(err)

	decoder := json.NewDecoder(f)
	config := config{}
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
	var configFile = flag.String("config", defaultConfig, "path to a JSON config file")
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
