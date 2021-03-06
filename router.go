package bellbox

import (
	//"fmt"
	"github.com/gin-gonic/gin"
)

func Route(router *gin.Engine) {
	// users
	router.POST("/user/new", HandleNewUser)
	router.POST("/user/login", HandleExistingUser)

	// user bells
	router.POST("/bell/new", HandleUserAuth(HandleNewBell))
	router.GET("/bell/map", HandleUserAuth(HandleMapBells))
	router.POST("/bell/delete", HandleUserAuth(HandleDeleteBell))

	// send
	router.POST("/send", HandleSenderAuth(HandleSend))

	// sender auth
	router.POST("/send/request", HandleSendRequest)
	router.GET("/send/map", HandleUserAuth(HandleMapAuthorizations))
	router.POST("/send/accept", HandleUserAuth(HandleSendAccept))
	router.POST("/send/deny", HandleUserAuth(HandleSendDeny))

	// messages
}
