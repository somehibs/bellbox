package bellbox

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleNewUser(c *gin.Context) {
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	if user.User == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required user or password"})
		return
	}
	// get a database, try add this person to it
	if UserExists(user.User) {
		c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
		return
	}
	var db = GetConfig().Db.GetDb()
	db.Create(&user)
	token := NewToken(user.User)
	ReplyToken(token, c)
}

func UserExists(user string) bool {
	fmt.Printf("Checking if user %s exists\n", user)
	var db = GetConfig().Db.GetDb()
	username := User{}
	db.Where("\"user\" = ?", user).Find(&username)
	if username.Password != "" {
		return true
	}
	return false
}

func NewToken(user string) string {
	token := GenToken()
	db.Create(&UserToken{User: user, Token: token, Timestamp: time.Now()})
	return token
}

func ReplyToken(token string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": token})
}

type GinHandler func(*gin.Context)

func HandleUserAuth(handler GinHandler) func(*gin.Context) {
	// Check if the user is permitted to continue
	return func(c *gin.Context) {
		a := c.Request.Header.Get("Authorization")
		if a == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth header not present"})
			return
		}
		//fmt.Println(a)
		db := GetConfig().Db.GetDb()
		token := UserToken{}
		db.Where("token = ?", a).Find(&token)
		if token.User == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "auth rejected"})
			return
		}
		fmt.Printf("Found user: %+v\n", token)
		c.Request.Header.Set("UserId", token.User)
		handler(c)
	}
}

func HandleExistingUser(c *gin.Context) {
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	if user.User == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields"})
		return
	}
	if !UserExists(user.User) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not exist or not match"})
		return
	}
	// get a database, try add this person to it
	var db = GetConfig().Db.GetDb()
	uexist := User{}
	db.Where("\"user\" = ?", user.User).Find(&uexist)
	if uexist.Password == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "not exist or not match"})
		return
	}
	if user.Password != uexist.Password {
		c.JSON(http.StatusForbidden, gin.H{"error": "not exist or not match"})
		return
	}
	// OK, make token
	token := NewToken(user.User)
	ReplyToken(token, c)
}
