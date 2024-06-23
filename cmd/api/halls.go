package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
)

func (server *Server) listHalls(ctx *gin.Context) {

	halls, err := server.store.Hall.ListHalls(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, halls)
}
func (server *Server) SearchHall(ctx *gin.Context) {
	name := ctx.Query("name")

	halls, err := server.store.Hall.SearchHall(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, halls)
}

func (server *Server) InsertHall(ctx *gin.Context) {
	var hall *models.Hall
	if err := ctx.ShouldBindJSON(&hall); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	if err := server.store.Hall.InsertHall(ctx, hall); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, hall)
}

func (server *Server) UpdateHall(ctx *gin.Context) {
	id := ctx.Param("id")
	var hall *models.Hall
	if err := ctx.ShouldBindJSON(&hall); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	modifiedCount, err := server.store.Hall.UpdateHall(ctx, id, *hall)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "hall not found",
		})
		return
	}

	ctx.JSON(http.StatusCreated, hall)
}

func (server *Server) DeleteHall(ctx *gin.Context) {
	id := ctx.Param("id")

	deletedCount, err := server.store.Hall.DeleteHall(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "book not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book has been deleted",
	})
}
