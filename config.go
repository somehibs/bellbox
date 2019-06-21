package bellbox

import (
	"encoding/json"
	"os"
	"io/ioutil"
)

type DbConfig struct {
	Host string
	Port string
	User string
	DbName string
	Password string
}

type PushConfig struct {
	// nothing yet
}

type Config struct {
	Db DbConfig
	Push PushConfig
}

var config Config

func GetConfig() Config {
	config := Config{}
	f, err := os.Open("config.json")
	if err != nil {
		panic("Could not open config.json. Please create from example config. e: " + err.Error())
	}
	fr, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Could not read anything from config.json. e: " + err.Error())
	}
	err = json.Unmarshal(fr, &config)
	if err != nil {
		panic("Could not understand JSON Config object from config.json. e: " + err.Error())
	}
	return config
}
