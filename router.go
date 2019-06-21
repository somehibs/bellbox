package bellbox

import (
	//"fmt"
	"github.com/gin-gonic/gin"
)

func Route(router *gin.Engine) {
	// TODO: HandleUserAuth(func(c*gin.Context))
	router.POST("/user/new", HandleNewUser)
	router.POST("/user/login", HandleExistingUser)
	//router.POST("
}
