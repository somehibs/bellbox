package bellbox

import (
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
	db.Model(&User{}).Create(&user)
}

func HandleExistingUser(c *gin.Context) {
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
}
