package bellbox

import (
	"fmt"
	"testing"
	"net/http"
	"strings"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

func TestAllMalformedJson(t *testing.T) {
	MalformedJsonCall(t, "user/new")
	MalformedJsonCall(t, "user/login")
}

func MalformedJsonCall(t *testing.T, path string) {
	buf := strings.NewReader("bitch")
	resp, _ := http.Post("http://localhost:8080/"+path, "application/json", buf)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status does not match expected")
	}
}

func TestGoodJson(t *testing.T) {
	u := User{User: "test", Password: "test"}
	buf, e:= json.Marshal(&u)
	if e != nil {
		t.Errorf("couldn't marshal user: %+v\n", e)
	}
	r, e := http.Post("http://localhost:8080/user/new", "application/json", bytes.NewReader(buf))
	if e != nil {
		t.Errorf("Expected success, received: %+v\n", e)
	}
	rb, e := ioutil.ReadAll(r.Body)
	fmt.Println(string(rb))
	r, e = http.Post("http://localhost:8080/user/login", "application/json", bytes.NewReader(buf))
	if e != nil {
		t.Errorf("Expected success, received: %+v\n", e)
	}
	rb, e = ioutil.ReadAll(r.Body)
	fmt.Println(string(rb))
	token := UserReply{}
	if e = json.Unmarshal(rb, &token); e != nil {
		t.Errorf("unmarshalling " + string(rb) + " failed")
	}
	fmt.Printf("%+v\n", token)
	badId := "wank"
	newBell := Bell{Id: badId, User: badId, Name: "test bell", Type: "ANDROID", Key: "fake", Enabled: false}
	if wb, e := json.Marshal(&newBell); e != nil {
		t.Errorf("marshalling %+v failed\n", newBell)
	} else {
		if r, e := Post(token.Token, "http://localhost:8080/bell/new", wb); e != nil {
			t.Errorf("failed to post %+v\n", e)
		} else {
			reply, _ := ioutil.ReadAll(r.Body)
			fmt.Printf("new bell response: %+v\n", string(reply))
		}
	}
	newBell = Bell{Id: badId, User: badId, Name: "web test bell", Type: "WEB", Key: "http://localhost:90210/bell", Enabled: false}
	wb, e := json.Marshal(&newBell)
	if e != nil {
		t.Errorf("marshalling %+v failed\n", newBell)
	}
	if r, e := Post(token.Token, "http://localhost:8080/bell/new", wb); e != nil || r.StatusCode != 200 {
		t.Errorf("creating bell failed %s", e)
	}

	r, e = Post(token.Token, "http://localhost:8080/bell/map", []byte{})
	reply, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("map: %+v\n", string(reply))

	sr := Bellringer{"test", "ringer test", "", false, 100}
	srr, e := json.Marshal(&sr)
	//fmt.Printf("send: %+v\n", string(srr))
	r, e = Post(token.Token, "http://localhost:8080/send/request", srr)
	reply, _ = ioutil.ReadAll(r.Body)
	bellringerRequest := UserReply{}
	json.Unmarshal(reply, &bellringerRequest)
	fmt.Printf("send response: %+v\n", bellringerRequest)

	r, e = Post(token.Token, "http://localhost:8080/send/map", []byte{})
	reply, _ = ioutil.ReadAll(r.Body)
	ringerList := []Bellringer{}
	e = json.Unmarshal(reply, &ringerList)
	fmt.Printf("auth map: %+v\n", ringerList)

	// test send
	msg := Message{"test", "Test title", "Test body", "normal"}
	SendTest(t, bellringerRequest.Token, msg, 403)

	// deny, test send, undeny
	buf, e = json.Marshal(ringerList[0])
	if e != nil {
		t.Errorf("error encoding json from response: %+v\n", e)
	}
	r, e = Post(token.Token, "http://localhost:8080/send/deny", buf)
	fmt.Printf("deny response: %d\n", r.StatusCode)

	SendTest(t, bellringerRequest.Token, msg, 403)

	r, e = Post(token.Token, "http://localhost:8080/send/accept", buf)
	fmt.Printf("accept response: %d\n", r.StatusCode)

	SendTest(t, bellringerRequest.Token, msg, 200)
}

func SendTest(t *testing.T, token string, msg Message, statusCode int) {
	buf, e := json.Marshal(msg)
	if e != nil {
		t.Errorf("Error marshalling: %+v\n", e)
		return
	}
	r, e := Post(token, "http://localhost:8080/send", buf)
	reply, _  := ioutil.ReadAll(r.Body)
	fmt.Printf("send code: %d reply: %+v\n", r.StatusCode, string(reply))
	if r.StatusCode != statusCode {
		t.Errorf("Expected status code not present (%d vs %d)", statusCode, r.StatusCode)
	}
}

func Post(token string, url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	c := http.Client{}
	r, e := c.Do(req)
	return r, e
}

