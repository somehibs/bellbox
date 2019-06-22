package bellbox

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleSendRequest(c *gin.Context) {
	ringer := Bellringer{}
	if err := c.ShouldBindJSON(&ringer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	// get a database, try add this person to it
	var db = GetConfig().Db.GetDb()
	find := Bellringer{Target: ringer.Target, Name: ringer.Name}
	db.Find(&find)
	if find.Token != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "sender already exists with same name", "detail": fmt.Sprintf("%+v\n", find)})
		return
	}
	ringer.RequestState = 0
	ringer.Token = GenToken()
	db.Create(&ringer)
	ReplyToken(ringer.Token, c)
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
		token := Bellringer{Token: a}
		db.Find(&token)
		if token.Name == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth rejected"})
			return
		} else if token.RequestState == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth pending"})
			return
		} else if token.RequestState == 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth rejected"})
			return
		}
		c.Request.Header.Set("Target", token.Target)
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
	sendMsgImpl(msg)
}

func sendMsgImpl(msg Message) {
	fmt.Printf("Sending message %+v\n", msg)
	var db = GetConfig().Db.GetDb()
	db.Create(&msg)
	// send to target bells
	bells := []Bell{}
	db.Where("\"user\" = ?", msg.Target).Find(&bells)
	for _, bell := range bells {
		fmt.Printf("Bell: %+v\n", bell)
		if bell.Type == "ANDROID" {
			PushAndroid(bell.Key, msg)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "target does not match id"})
		return
	}
	db := GetConfig().Db.GetDb()
	db.Find(&ringer)
	if ringer.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ringer does not exist"})
		return
	}
	if enable {
		ringer.RequestState = 1
	} else {
		ringer.RequestState = 2
	}
	db.Table("bellringers").Where("target = ? AND name = ?", ringer.Target, ringer.Name).Update("request_state", ringer.RequestState)
	rr := Bellringer{}
	db.Find(&rr)
	if rr.RequestState != ringer.RequestState {
		c.JSON(http.StatusConflict, gin.H{"error": "ringer not in expected state after transaction"})
		return
	}
	fmt.Printf("post update ringer: %+v\n", rr)
}

func HandleMapAuthorizations(c *gin.Context) {
	user := c.Request.Header.Get("UserId")
	ringerSlice := make([]Bellringer, 0)
	db.Where("target = ?", user).Find(&ringerSlice)
	fmt.Printf("ringers: %+v\n", ringerSlice)
	c.JSON(http.StatusOK, &ringerSlice)
}