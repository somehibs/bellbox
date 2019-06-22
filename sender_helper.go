package bellbox

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type SenderAuth struct {
	Name string
	Keys map[string]string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadSenderAuth() SenderAuth {
	auth := SenderAuth{"",map[string]string{}}
	f, err := os.Open("bellbox_auth.json")
	if err != nil {
		fmt.Printf("file not exist?: %s\n", err)
		return auth
	}
	read, err := ioutil.ReadAll(f)
	check(err)
	err = json.Unmarshal(read, &auth)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling: %s\n", err))
	}
	return auth
}

func UpdateSenderAuth(auth SenderAuth) {
	os.Remove("bellbox_auth.json")
	f, err := os.Create("bellbox_auth.json")
	check(err)
	fmt.Printf("Writing auth %+v to file\n", auth)
	buf, err := json.Marshal(&auth)
	fmt.Printf("Writing to file: %s\n", buf)
	check(err)
	_, err = f.Write(buf)
	check(err)
	f.Close()
}

var auth SenderAuth
var server string

func StartSender(_name string, _server string) {
	auth = LoadSenderAuth()
	if auth.Name != "" && auth.Name != _name {
		panic("config name doesn't match provided name - auth will no longer be valid")
	}
	server = _server
	auth.Name = _name
}

func checkSender() {
	if server == "" {
		panic("Server url not set - call StartSender first")
	}
	if auth.Name == "" {
		panic("Sender not set - call StartSender first")
	}
}

func SingleTarget(targetName string) (string, error) {
	checkSender()
	Log("Creating a single target.")
	// connect to server, request target
	if auth.Keys[targetName] != "" {
		Log("Found existing key.")
		return auth.Keys[targetName], nil
	}
	bellringer := Bellringer{targetName, auth.Name, "", false, 100}
	reply := UserReply{}
	_, e := Post("", server+"/send/request", &bellringer, &reply)
	if e != nil {
		Log("Failed to request target.")
		return "", e
	}
	auth.Keys[targetName] = reply.Token
	UpdateSenderAuth(auth)
	Log("Request complete, token: " + reply.Token)
	return reply.Token, nil
}

func Send(target, targetToken, title, message string) error {
	msg := Message{Title: title, Message: message, Target: target}
	_, err := Post(targetToken, server+"/send", &msg, nil)
	return err
}
