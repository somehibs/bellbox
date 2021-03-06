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
	Fcm string
}

type Config struct {
	Db DbConfig
	Push PushConfig
}

var config Config
var loaded bool

func GetConfig() Config {
	if loaded {
		return config
	}
	config = Config{
		Db: DbConfig {
			Host: "localhost",
			Port: "5432",
			DbName: "bellbox",
		},
		Push: PushConfig{},
	}
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
	loaded = true
	return config
}
