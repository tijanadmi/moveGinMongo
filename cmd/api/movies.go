package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
)

func (server *Server) searchMovies(ctx *gin.Context) {

	movieId := ctx.Param("id")
	// Provera da li je movieId "%", ako jeste, postavlja se na prazan string
	// if movieId == "%" {
	// 	movieId = ""
	// }
	movies, err := server.store.Movie.SearchMovies(ctx, movieId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

func (server *Server) listMovies(ctx *gin.Context) {

	movies, err := server.store.Movie.SearchMovies(ctx, "0")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

func (server *Server) InsertMovie(ctx *gin.Context) {
	var movie *models.Movie
	if err := ctx.ShouldBindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	if err := server.store.Movie.AddMovie(ctx, movie); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, movie)
}

func (server *Server) UpdateMovie(ctx *gin.Context) {
	id := ctx.Param("id")
	var movie *models.Movie
	if err := ctx.ShouldBindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	modifiedCount, err := server.store.Movie.UpdateMovie(ctx, id, *movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (server *Server) DeleteMovie(ctx *gin.Context) {
	id := ctx.Param("id")

	deletedCount, err := server.store.Movie.DeleteMovie(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Movie has been deleted",
	})
}
