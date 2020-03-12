package bellbox

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

func Log(s string) {
	fmt.Printf("bellbox: %s\n", s)
}

func Post(token string, url string, body interface{}, reply interface{}) (*http.Response, error) {
	bbody, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bbody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	c := http.Client{}
	r, e := c.Do(req)
	if e != nil {
		return nil, e
	}
	read , _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != 200 {
		fmt.Println("Status code did not match. Cannot continue.")
		fmt.Println("Status code did not match. Cannot continue.")
		panic(fmt.Sprintf("status code is %d body is %s\n", r.StatusCode, read))
		fmt.Println("Status code did not match. Cannot continue.")
		fmt.Println("Status code did not match. Cannot continue.")
	}
	e = json.Unmarshal(read, &reply)
	return r, e
}
