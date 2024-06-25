package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
)

func (server *Server) GetRepertoire(ctx *gin.Context) {

	id := ctx.Param("id")
	repertoire, err := server.store.Repertoire.GetRepertoire(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, repertoire)
}

func (server *Server) GetAllRepertoireForMovie(ctx *gin.Context) {

	movieId := ctx.Query("movie_id")

	repertoires, err := server.store.Repertoire.GetAllRepertoireForMovie(ctx, movieId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, repertoires)
}

func (server *Server) ListRepertoires(ctx *gin.Context) {

	repertoires, err := server.store.Repertoire.ListRepertoires(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, repertoires)
}

func (server *Server) AddRepertoire(ctx *gin.Context) {
	var repertoire *models.Repertoire
	if err := ctx.ShouldBindJSON(&repertoire); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	fmt.Println("add handler", repertoire.NumOfResTickets)
	if err := server.store.Repertoire.AddRepertoire(ctx, repertoire); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, repertoire)
}

func (server *Server) UpdateRepertoire(ctx *gin.Context) {
	id := ctx.Param("id")
	var repertoire *models.Repertoire
	if err := ctx.ShouldBindJSON(&repertoire); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}
	modifiedCount, err := server.store.Repertoire.UpdateRepertoire(ctx, id, *repertoire)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "repertoire not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, repertoire)
}

func (server *Server) DeleteRepertoire(ctx *gin.Context) {
	id := ctx.Param("id")

	deletedCount, err := server.store.Repertoire.DeleteRepertoire(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "repertoires not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d repertoire has been deleted", deletedCount),
	})
}

func (server *Server) DeleteRepertoireForMovie(ctx *gin.Context) {
	movieId := ctx.Query("movie_id")

	deletedCount, err := server.store.Repertoire.DeleteRepertoireForMovie(ctx, movieId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "repertoires not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d repertoire has been deleted", deletedCount),
	})
}
