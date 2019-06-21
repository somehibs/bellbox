package main

import (
	"git.circuitco.de/self/bellbox"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	bellbox.Route(router)
	router.Run()
}
