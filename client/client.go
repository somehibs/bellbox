package main

import (
	"fmt"
	"net/http"
	//"strings"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"git.circuitco.de/self/bellbox"
	"flag"
)

func Errorf(str string, kwarg... interface{}) {
	fmt.Printf(str, kwarg)
	panic("")
}

func LoginUser(user bellbox.User, create bool) string {
	u := user
	buf, e:= json.Marshal(&u)
	if e != nil {
		Errorf("couldn't marshal user: %+v\n", e)
	}
	path := "login"
	if create {
		path = "new"
	}
	r, e := http.Post("http://localhost:8080/user/"+path, "application/json", bytes.NewReader(buf))
	if e != nil {
		Errorf("Expected success, received: %+v\n", e)
	}
	rb, e := ioutil.ReadAll(r.Body)
	fmt.Println(string(rb))
	token := bellbox.UserReply{}
	if e = json.Unmarshal(rb, &token); e != nil {
		Errorf("unmarshalling " + string(rb) + " failed")
	}
	fmt.Printf("%+v\n", token)
	return token.Token
}

func ListBell(token string) []bellbox.Bell {
	r, e := Post(token, "http://localhost:8080/bell/map", []byte{})
	reply, _ := ioutil.ReadAll(r.Body)
	bells := []bellbox.Bell{}
	fmt.Printf("map str: %s\n", reply)
	e = json.Unmarshal(reply, &bells)
	fmt.Printf("map: %+v\n", bells)
	if e != nil {
		panic("could not list bells")
	}
	return bells
}

func ListAuthMap(token string) []bellbox.Bellringer {
	r, e := Post(token, "http://localhost:8080/send/map", []byte{})
	if e != nil {
		panic(e.Error())
	}
	reply, _ := ioutil.ReadAll(r.Body)
	ringerList := []bellbox.Bellringer{}
	e = json.Unmarshal(reply, &ringerList)
	if e != nil {
		panic(e.Error())
	}
	fmt.Printf("auth map: %+v\n", ringerList)
	return ringerList
}

//func TestSend() {
//	// test send
//	msg := bellbox.Message{"test", "Test title", "Test body", "normal"}
//	SendTest(t, bellringerRequest.Token, msg, 403)
//
//	// deny, test send, undeny
//	buf, e = json.Marshal(ringerList[0])
//	if e != nil {
//		t.Errorf("error encoding json from response: %+v\n", e)
//	}
//	r, e = Post(token.Token, "http://localhost:8080/send/deny", buf)
//	fmt.Printf("deny response: %d\n", r.StatusCode)
//
//	SendTest(t, bellringerRequest.Token, msg, 403)
//
//	r, e = Post(token.Token, "http://localhost:8080/send/accept", buf)
//	fmt.Printf("accept response: %d\n", r.StatusCode)
//
//	SendTest(t, bellringerRequest.Token, msg, 200)
//}

func Post(token string, url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	c := http.Client{}
	r, e := c.Do(req)
	return r, e
}

func main() {
	uflag := flag.String("user", "", "Username for login/new operations")
	pflag := flag.String("pass", "", "Password for login/new operations")
	mode := flag.String("mode", "", "Mode flag. Needs to be a) new b) login c) bells d) auths e) accept f) deny")
	flag.Parse()
	fmt.Printf("Found flags: u: %s p: %s\n", *uflag, *pflag)
}
