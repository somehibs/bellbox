package bellbox

import (
	//"fmt"
	"github.com/gin-gonic/gin"
)

func Route(router *gin.Engine) {
	// TODO: HandleUserAuth(func(c*gin.Context))
	// users
	router.POST("/user/new", HandleNewUser)
	router.POST("/user/login", HandleExistingUser)

	router.POST("/bell/new", HandleUserAuth(HandleNewBell))
	router.POST("/bell/map", HandleUserAuth(HandleMapBells))
}
