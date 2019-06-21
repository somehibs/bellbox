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

func HandleSend(c *gin.Context) {
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	if user.User == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields"})
	}
	// get a database, try add this person to it
	var db = GetConfig().Db.GetDb()
	db.Find(&user)
	if user.User == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "not exist or not match"})
		return
	}
	// OK, make token
	token := NewToken(user.User)
	ReplyToken(token, c)
}

func HandleSendAuthorizations(c *gin.Context) {

}

func HandleSendAccept(c *gin.Context) {
}

func HandleSendDeny(c *gin.Context) {
}
