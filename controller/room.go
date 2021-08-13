package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shreyngd/booker/db"
	"github.com/shreyngd/booker/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Controller) AddRoom(ctx *gin.Context) {
	var room models.InterviewRoom
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data": err,
		})
		return
	}

	room.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	room.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	room.ID = primitive.NewObjectID()
	room.Active = true

	log.Println(room, "room")

	if err := db.CreateRoom(&room); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"data": err,
		})
		return
	}
	gb := models.GetInstanceGlobal()
	gb.Channel <- *room.Name
	ctx.JSON(http.StatusCreated, gin.H{
		"data": "Room Created Successfully",
	})

}
