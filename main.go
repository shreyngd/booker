package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shreyngd/booker/controller"
)


func main() {
	r := gin.Default()
	r.GET("/test", handleFunc)
	c := controller.NewController()

	v1 := r.Group("/api/v1")
	{
		books:= v1.Group("/books")
		{
			books.GET("/",c.GetBooks)
			books.POST("/",c.PutBooks)
		}

	}
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}

func handleFunc(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok!!!",
	})
}