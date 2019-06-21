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
}
