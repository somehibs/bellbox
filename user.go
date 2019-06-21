package bellbox

import (
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleNewUser(c *gin.Context) {
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	// get a database, try add this person to it
	var db = GetConfig().Db.GetDb()
	username := User{User: user.User}
	db.Find(&username)
	if username.Password != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
		return
	}
	db.Create(&user)
	token := NewToken(user.User)
	ReplyToken(token, c)
}

func NewToken(user string) string {
	token := GenToken()
	db.Create(&UserToken{User:user, Token:token, Timestamp: time.Now()})
	return token
}

func ReplyToken(token string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func HandleExistingUser(c *gin.Context) {
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
