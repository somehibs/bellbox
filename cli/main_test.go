package main

import (
	"fmt"
	"testing"
	"net/http"
	"strings"
	"bytes"
	"io/ioutil"
	"encoding/json"

	"git.circuitco.de/self/bellbox/api"
)


func TestMalformedJson(t *testing.T) {
	buf := strings.NewReader("bitch")
	_, e := http.Post("http://localhost:8080/user/new", "application/json", buf)
	if e == nil {
		t.Errorf("Expected error, received: %+v\n", e)
	}
}

func TestGoodJson(t *testing.T) {
	u := api.User{"test", "test", false}
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
}
