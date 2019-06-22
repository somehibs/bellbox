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
}

func HandleMapBells(c *gin.Context) {
	bellSlice := make([]Bell, 0)
	db.Find(&bellSlice)
	fmt.Printf("bells: %+v\n", bellSlice)
	c.JSON(http.StatusOK, &bellSlice)
}
