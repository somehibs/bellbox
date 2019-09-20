package bellbox

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type AndroidPush struct {
	Sender    string    `json:"sender"`
	Body      string    `json:"body"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"time"`
	Url       string    `json:"url"`
}

func PushAndroid(token string, msg Message) {
	if GetConfig().Push.Fcm == "" {
		panic("Cannot continue - FCM key missing for Android push")
	}
	req := `{"to":"` + token + `", "data":`
	push := AndroidPush{msg.Sender, msg.Message, msg.Title, msg.Timestamp, msg.Url}
	marshalledJson, err := json.Marshal(&push)
	req += string(marshalledJson) + "}"
	r, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", strings.NewReader(req))
	if err != nil {
		panic("cannot construct request")
	}
	r.Header.Set("Authorization", "key="+GetConfig().Push.Fcm)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Connection", "close")
	c := http.Client{Timeout: time.Second * 5}
	resp, err := c.Do(r)
	if err != nil {
		panic("could not do request")
	}
	resp.Body.Close()
}
