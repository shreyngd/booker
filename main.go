package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shreyngd/booker/controller"
	"github.com/shreyngd/booker/middleware"
	slotserver "github.com/shreyngd/booker/slotServer"
)

func main() {
	r := gin.Default()

	c := controller.NewController()

	v1 := r.Group("/api/v1")

	wsServer := slotserver.NewWebsocketServer()
	go wsServer.Run()

	v1.GET("/ws", func(c *gin.Context) {
		slotserver.ServeWs(wsServer, c.Writer, c.Request)
	})
	v1.GET("/test", handleFunc)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/signup", c.Signup)
		auth.GET("/google", c.LoginGoogle)
		auth.GET("/google/callback", c.CallbackGoogle)

	}
	v1.Use(middleware.Authentication())
	books := v1.Group("/books")
	{
		books.GET("/", c.GetBooks)
		books.POST("/", c.PutBooks)
	}
	rooms := v1.Group("/room")
	{
		rooms.POST("/", c.AddRoom)
	}

	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}

func handleFunc(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok!!!",
	})
}
