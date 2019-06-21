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
	newBell := Bell{Id: badId, User: badId, Name: "test bell", Type: "WEB", Key: "http://wank", Enabled: false}
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
	r, e = Post(token.Token, "http://localhost:8080/bell/map", []byte{})
	reply, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("map: %+v\n", string(reply))
}

func Post(token string, url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	c := http.Client{}
	return c.Do(req)
}

