package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shreyngd/booker/db"
	"github.com/shreyngd/booker/models"
)

// Get all Books controller
func (c *Controller) GetBooks(ctx  *gin.Context){
	books, err := db.GetBooks()
	

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": err,
		})
		return
	}
	if len(books) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": -3,
			"data": books,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"data": books,
	})
}

//Insert Many Books in the db
func (c *Controller) PutBooks(ctx *gin.Context){
	var books []models.Book
	if err := ctx.ShouldBindJSON(&books); err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error": err,
		})
		return
	}
	if len(books) == 0{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "books are empty",
			},
		})	
	}
	count, err := db.PutBooks(books)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": -1,
			"data": gin.H{
				"error": err,
			},
		})		
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"data": gin.H{
			"count": count,
		},
	})
}