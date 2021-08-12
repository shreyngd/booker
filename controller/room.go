package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shreyngd/booker/models"
)

func (c *Controller) AddRoom(ctx *gin.Context) {
	var room models.InterviewRoom
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data": err,
		})
		return
	}

}
