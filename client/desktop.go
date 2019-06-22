package main

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"os"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"git.circuitco.de/self/bellbox"
)

type ClientConfig struct {
	// only needed for belltoken setup, will be omitted when saving back to json
	CreateUser bool `json:",omitempty"`
	User string `json:",omitempty"`
	Pass string `json:",omiempty"`
	// bell specific token for use at some point
	BellToken string
	// server-accessible host
	BellHost string
	// used to identify the bell to the user
	Name string
	// server
	Url string
}

func ReadConfig() ClientConfig {
	config := ClientConfig{}
	f, err := os.Open("bell_config.json")
	if err != nil {
		panic("Could not open bell_config.json")
	}
	c, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Could not read bell_config.json")
	}
	if e := json.Unmarshal(c, &config); e != nil {
		panic("Could not de-json bell_config.json ("+e.Error()+")")
	}
	return config
}

func SaveConfig(config ClientConfig) {
	f, err := os.Open("bell_config.json")
	if err != nil {
		panic("Could not open bell_config.json")
	}
	c, err := json.Marshal(config)
	if err != nil {
		panic("Could not marshal into json: "+ err.Error())
	}
	if _, err := f.Write(c); err != nil {
		panic(err.Error())
	}
	f.Close()
}
func CreateBell(config ClientConfig, token string) (string, error) {
	req, e := json.Marshal(bellbox.Bell{Name:config.Name, Key: config.BellHost, Type: "WEB"})
	if e != nil {
		return "", e
	}
	re, e := http.NewRequest("POST", config.Url+"bell/new", bytes.NewReader(req))
	re.Header.Set("Authorization", token)
	cl := http.Client{}
	r, e := cl.Do(re)
	if e != nil {
		return "", e
	}
	read, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("registration: %d resp: %s\n", r.StatusCode, read)
	tok := bellbox.UserToken{}
	e = json.Unmarshal(read, &tok)
	if e != nil {
		return "", e
	}
	return tok.Token, e
}


func CreateUser(root, user, pass string) (string, error) {
	req, e := json.Marshal(bellbox.User{User: user, Password: pass})
	if e != nil {
		return "", e
	}
	r, e := http.Post(root+"user/new", "application/json", bytes.NewReader(req))
	if e != nil {
		return "", e
	}
	read, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("registration: %d resp: %s\n", r.StatusCode, read)
	tok := bellbox.UserToken{}
	e = json.Unmarshal(read, &tok)
	if e != nil {
		return "", e
	}
	return tok.Token, nil
}

func ListenForBells() {
	// Open an http server that can handle that request
	r := gin.Default()
	fmt.Printf("bells: %+v\n", r)
	r.POST("/bell", func(c *gin.Context) {
		fmt.Printf("req: %+v\n", c.Request)
	})
	r.Run()
}

func main() {
	// Read configuration from disk
	config := ReadConfig()
	fmt.Printf("Config: %+v\n", config)
	if config.BellToken == "" {
		fmt.Println("No bell token found. Attempting creation.")
		if false {//config.User == "" {//|| config.Pass == "" {
			panic("User credentials not present, cannot create bell.")
		}
		if config.Url == "" {
			panic("bellbox host not found")
		}
		if config.BellHost == "" {
			panic("bellbox local host not found")
		}
		token := "Bi93gKjvmh4vPziMcXYwOPeAOPqikRiNzYRUBwOYD7MWpop3y5sLmsvLp2nRoPke"
		var err error
		if config.CreateUser {
			fmt.Println("User registration requested.")
			if token, err = CreateUser(config.Url, config.User, config.Pass); err != nil {
				panic("Could not create user: " + err.Error())
			}
		}
		if token == "" {
			panic("Could not obtain token, bell cannot be registered, halting")
		} else {
			// We have a user token, create the bell
			CreateBell(config, token)
			config.User = ""
			config.Pass = ""
			config.CreateUser = false
			config.BellToken = token
			SaveConfig(config)
		}
	} else {
		ListenForBells()
	}
}
