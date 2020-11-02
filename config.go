package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/google/uuid"
)

type Config struct {
	UserID         string `json:"user_id"`
	Name           string
	Email          string
	ServiceURL     string         `json:"service_url"`
	TimeCategories []TimeCategory `json:"time_categories"`
}

type TimeCategory struct {
	ID     string
	Name   string
	Color  string
	Active bool
}

func getConfig() Config {
	// Get the current user
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Unable to get the active user!")
		panic(err)
	}
	configFile := usr.HomeDir + "/.bt"
	// Check for the existence of a .bt config file in the user's home dir
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create a default configuration
		defaultConfigFor(configFile)
		fmt.Println("\nYou have a new default config file at: " + configFile)
		fmt.Print("\nPlease edit the file to match your desired configuration.\n\n")
		os.Exit(0)
	}
	// Read the config file
	fileContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Unable to read the config file at: " + configFile)
		panic(err)
	}
	// Parse the config file JSON
	conf := Config{}
	err = json.Unmarshal([]byte(fileContents), &conf)
	if err != nil {
		fmt.Println("Unable to parse the config file at: " + configFile)
		panic(err)
	}
	// TODO - Validate the data contents
	return conf
}

func defaultConfigFor(configFile string) {
	// Read the default config file
	// TODO what happens to this asset file when it's packaged up for usage?
	fileContents, err := ioutil.ReadFile("./assets/default.cfg.json")
	if err != nil {
		fmt.Println("Unable to read the default config file asset.")
		panic(err)
	}
	// Parse the default config file JSON
	userConf := Config{}
	err = json.Unmarshal([]byte(fileContents), &userConf)
	if err != nil {
		fmt.Println("Unable to parse the default config file asset.")
		panic(err)
	}
	// Replace the user with a new UUID
	userConf.UserID = uuid.New().String()
	for i, category := range userConf.TimeCategories {
		category.ID = uuid.New().String()
		userConf.TimeCategories[i] = category
	}

	// Write the user's new config file
	conf, err := json.MarshalIndent(userConf, "", " ")
	err = ioutil.WriteFile(configFile, conf, 0644)
}