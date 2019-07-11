package bellbox

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type SenderAuth struct {
	// Map of sender names to sender single token auths
	SingleSenders map[string]map[string]string
	// bellbox server ip
	Server string
	CurrentSender string `json:""`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadSenderAuth() SenderAuth {
	auth := SenderAuth{map[string]map[string]string{}, "", ""}
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

var server string

func StartSender(_name string, _server string) SenderAuth {
	auth := LoadSenderAuth()
	// Check if we've already got a sender by this name
	auth.Server = _server
	auth.CurrentSender = _name
	if auth.SingleSenders[_name] == nil {
		auth.SingleSenders[_name] = map[string]string{}
	}
	return auth
}

func (auth *SenderAuth) checkServer() {
	if auth.Server == "" {
		panic("Server url not set - call StartSender first")
	}
}

func (auth *SenderAuth) SingleTarget(targetName string) (string, error) {
	auth.checkServer()
	Log("Creating a single target.")
	// connect to server, request target
	if auth.SingleSenders[auth.CurrentSender][targetName] != "" {
		Log("Found existing key.")
		return auth.SingleSenders[auth.CurrentSender][targetName], nil
	}
	bellringer := Bellringer{targetName, auth.CurrentSender, "", false, 100}
	reply := UserReply{}
	_, e := Post("", server+"/send/request", &bellringer, &reply)
	if e != nil {
		Log("Failed to request target.")
		return "", e
	}
	auth.SingleSenders[auth.CurrentSender][targetName] = reply.Token
	UpdateSenderAuth(*auth)
	Log("Request complete, token: " + reply.Token)
	return reply.Token, nil
}

func Send(target, targetToken, title, message string) error {
	msg := Message{Title: title, Message: message, Target: target}
	_, err := Post(targetToken, server+"/send", &msg, nil)
	return err
}
