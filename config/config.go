package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Conf is the global config file
var (
	Conf Config
)

// Load is used to load the config toml file
func Load() {
	raw, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	if _, err = toml.Decode(string(raw), &Conf); err != nil {
		panic(err)
	}

}
