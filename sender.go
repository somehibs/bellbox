package bellbox

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var SystemSender = "System"

func HandleSendRequest(c *gin.Context) {
	ringer := Bellringer{}
	if err := c.ShouldBindJSON(&ringer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	if ringer.Target == "" || ringer.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required target or name"})
		return
	}
	if !UserExists(ringer.Target) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target does not exist"})
		return
	}
	if strings.Compare(SystemSender, ringer.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must not be " + SystemSender})
		return
	}
	// get a database, try add this person to it
	var db = GetConfig().Db.GetDb()
	find := Bellringer{}
	db.Where("target = ?", ringer.Target).Where("name = ?", ringer.Name).Find(&find)
	if find.Token != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "sender already exists with same name", "detail": fmt.Sprintf("%+v\n", find)})
		return
	}
	ringer.RequestState = 0
	ringer.Token = GenToken()
	db.Create(&ringer)
	systemMessage(ringer.Target, "New bellringer request", fmt.Sprintf("%s wants to send you notifications", ringer.Name))
	ReplyToken(ringer.Token, c)
}

func systemMessage(target, title, msg string) {
	// send message
	sendMsgImpl(Message{Target: target, Title: title, Message: msg, Sender: SystemSender})
}

func HandleSenderAuth(handler GinHandler) func(*gin.Context) {
	// Check if the bellringer is permitted to continue
	return func(c *gin.Context) {
		a := c.Request.Header.Get("Authorization")
		if a == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth header not present"})
			return
		}
		db := GetConfig().Db.GetDb()
		token := Bellringer{}
		db.Where("token = ?", a).Find(&token)
		if token.Name == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth does not exist"})
			return
		} else if token.RequestState == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth pending"})
			return
		} else if token.RequestState == 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth rejected"})
			return
		}
		c.Request.Header.Set("Target", token.Target)
		c.Request.Header.Set("Sender", token.Name)
		if token.Urgent {
			c.Request.Header.Set("Urgent", "true")
		}
		handler(c)
	}
}

var ValidPriorities = []string{"normal"}

func ValidPriority(priority string) bool {
	for _, valid := range ValidPriorities {
		if priority == valid {
			return true
		}
	}
	return false
}

func HandleSend(c *gin.Context) {
	msg := Message{}
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	if msg.Target == "" || msg.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields"})
		return
	}
	if msg.Priority == "" {
		msg.Priority = "normal"
	}
	if !ValidPriority(msg.Priority) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid priority"})
		return
	}
	if msg.Target != c.Request.Header.Get("Target") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message target does not match authentication"})
		return
	}
	msg.Sender = c.Request.Header.Get("Sender")
	sendMsgImpl(msg)
	c.JSON(http.StatusOK, gin.H{})
}

func sendMsgImpl(msg Message) {
	//fmt.Printf("Sending message %+v\n", msg)
	var db = GetConfig().Db.GetDb()
	msg.Timestamp = time.Now()
	db.Create(&msg)
	// send to target bells
	bells := []Bell{}
	db.Where("\"user\" = ?", msg.Target).Find(&bells)
	for _, bell := range bells {
		if bell.Type == "ANDROID" {
			go PushAndroid(bell.Key, msg)
		} else if bell.Type == "WEB" {
			go PushWeb(bell.Key, msg)
		}
	}
}

func HandleSendAccept(c *gin.Context) {
	HandleSendChange(c, true)
}

func HandleSendDeny(c *gin.Context) {
	HandleSendChange(c, false)
}

func HandleSendChange(c *gin.Context, enable bool) {
	ringer := Bellringer{}
	if err := c.ShouldBindJSON(&ringer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	target := c.Request.Header.Get("UserId")
	if ringer.Target != target {
		ringer.Target = target
	}
	db := GetConfig().Db.GetDb()
	existingRinger := Bellringer{}
	db.Where("target = ?", ringer.Target).Where("name = ?", ringer.Name).Find(&existingRinger)
	if existingRinger.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ringer does not exist"})
		return
	}
	actionName := "allowed"
	if enable {
		existingRinger.RequestState = 1
	} else {
		existingRinger.RequestState = 2
		actionName = "denied"
	}
	db.Table("bellringers").Where("target = ? AND name = ?", existingRinger.Target, existingRinger.Name).Update("request_state", existingRinger.RequestState)
	rr := Bellringer{}
	db.Where("target = ?", ringer.Target).Where("name = ?", ringer.Name).Find(&rr)
	if rr.RequestState != existingRinger.RequestState {
		c.JSON(http.StatusConflict, gin.H{"error": "ringer not in expected state after transaction"})
		return
	}
	systemMessage(ringer.Target, fmt.Sprintf("Bellringer %s", actionName), fmt.Sprintf("Bellringer %s was %s", ringer.Name, actionName))
	//fmt.Printf("post update ringer: %+v\n", rr)
	c.JSON(http.StatusOK, gin.H{})
}

func HandleMapAuthorizations(c *gin.Context) {
	user := c.Request.Header.Get("UserId")
	ringerSlice := make([]Bellringer, 0)
	db.Where("target = ?", user).Find(&ringerSlice)
	//fmt.Printf("ringers: %+v\n", ringerSlice)
	c.JSON(http.StatusOK, &ringerSlice)
}
