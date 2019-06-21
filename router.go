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

	router.POST("/send", HandleSenderAuth(HandleSend))
	router.POST("/send/request", HandleSendRequest)
	router.POST("/send/map", HandleUserAuth(HandleSendAuthorizations))
	router.POST("/send/accept", HandleUserAuth(HandleSendAccept))
	router.POST("/send/deny", HandleUserAuth(HandleSendDeny))
}
