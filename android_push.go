package bellbox

import (
	"fmt"
)
import "net/http"
import "strings"
import "io/ioutil"

func PushAndroid(token string, msg Message) {
	if GetConfig().Push.Fcm == "" {
		panic("Cannot continue - FCM key missing for Android push")
	}
	req := `{"to":"` + token + `", "data":{`
	req += `"title": "` + msg.Title + `"`
	if msg.Message != "" {
		req += `,"body": "` + msg.Message + `"`
	}
	req += `,"sender": "` + msg.Sender + `"`
	req += `}}`
	r, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", strings.NewReader(req))
	if err != nil {
		panic("cannot construct request")
	}
	r.Header.Set("Authorization", "key="+GetConfig().Push.Fcm)
	r.Header.Set("Content-Type", "application/json")
	c := http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		panic("could not do request")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp: %d body: %s", resp.StatusCode, body)
}
