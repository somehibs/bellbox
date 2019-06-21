package bellbox

import (
	"fmt"
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
	fmt.Println(db)
}

func HandleExistingUser(c *gin.Context) {
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed json"})
		return
	}
}
