package bellbox

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleNewBell(c *gin.Context) {
	authedUser := c.Request.Header.Get("UserId")
	bell := Bell{}
	if err := c.ShouldBindJSON(&bell); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	if bell.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing name"})
		return
	}
	// get a database, try add this person to it
	var db = GetConfig().Db.GetDb()
	bellNameCheck := Bell{}
	db.Where("name = ?", bell.Name).Where("\"user\" = ?", authedUser).Find(&bellNameCheck)
	if bellNameCheck.Id != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
		return
	}
	bell.Id = GenToken()
	bell.User = authedUser
	bell.Enabled = true
	db.Create(&bell)
	c.JSON(http.StatusOK, bell)
	systemMessage(authedUser, "Bell added", bell.Name+" has been added to your account.")
}

func HandleDeleteBell(c *gin.Context) {
	authedUser := c.Request.Header.Get("UserId")
	bell := Bell{}
	if err := c.ShouldBindJSON(&bell); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	if bell.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing key"})
		return
	}
	// get a database, try delete that bell
	var db = GetConfig().Db.GetDb()
	bellCheck := Bell{}
	db.Where("name = ?", bell.Name).Where("key = ?", bell.Key).Where("\"user\" = ?", authedUser).Find(&bellCheck)
	if bellCheck.Id == "" {
		c.JSON(http.StatusConflict, gin.H{"error": "bell does not exist"})
		return
	}
	systemMessage(authedUser, "Bell deleted", bell.Name+" has been deleted from your account.")
	db.Where("id = ?", bellCheck.Id).Delete(Bell{})
	c.JSON(http.StatusOK, gin.H{})
}

func HandleMapBells(c *gin.Context) {
	bellSlice := make([]Bell, 0)
	db.Find(&bellSlice)
	fmt.Printf("bells: %+v\n", bellSlice)
	c.JSON(http.StatusOK, &bellSlice)
}
