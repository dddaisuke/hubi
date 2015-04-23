package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	AccessToken string
}

func Parse(filename string) (Config, error) {
	var c Config
	jsonString, err := ioutil.ReadFile(filename)

	if err != nil {
		return c, err
	}
	err = json.Unmarshal(jsonString, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
