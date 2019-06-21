package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.POST("/user/new", func(c *gin.Context) {
		c.JSON(200, gin.H{"automatic": "reply"})
	})
	r.Run()
}
