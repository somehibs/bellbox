package bellbox

import (
	"fmt"
)
import "net/http"
import "strings"
import "io/ioutil"

type AndroidPush struct {
	sender string
	body string
	title string
	timestamp time.Time
}

func PushAndroid(token string, msg Message) {
	if GetConfig().Push.Fcm == "" {
		panic("Cannot continue - FCM key missing for Android push")
	}
	req := `{"to":"` + token + `", "data":`
	push := AndroidPush{msg.Sender, msg.Message, msg.Title, msg.Timestamp}
	marshalledJson, err := json.Marshal(&push)
	req += marshalledJson + "}"
	fmt.Println("ABout to send: " + req)
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
