package routes

import (
	"server/handlers"
	"server/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	wsManger := websocket.NewManager()
	r.GET("/health", handlers.Health)
	r.GET("/ws", wsManger.ServeWS)
}
