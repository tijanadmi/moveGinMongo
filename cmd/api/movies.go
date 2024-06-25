package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
