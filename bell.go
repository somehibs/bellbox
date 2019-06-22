package bellbox

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleNewBell(c *gin.Context) {
	bell := Bell{}
	if err := c.ShouldBindJSON(&bell); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
	// get a database, try add this person to it
	var db = GetConfig().Db.GetDb()
	bellNameCheck := Bell{}
	db.Where("name = ?", bell.Name).Where("user = ?", bell.User).Find(&bellNameCheck)
	if bellNameCheck.Id != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
		return
	}
	bell.Id = GenToken()
	bell.User = c.Request.Header.Get("UserId")
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
