package bellbox

import (
	"fmt"
	"net/http"
	"bytes"
	//"strings"
	"encoding/json"
	"io/ioutil"
)

func PushWeb(token string, msg Message) error {
	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		panic("cannot marshal json")
		return err
	}
	r, err := http.NewRequest("POST", token, bytes.NewReader(encodedMsg))
	if err != nil {
		panic("cannot construct request")
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	c := http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp: %d body: %s", resp.StatusCode, body)
	return nil
}
