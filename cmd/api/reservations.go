package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
	"github.com/tijanadmi/moveginmongo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reservationRequest struct {
	Username    string   `json:"username" binding:"required"`
	MovieID     string   `json:"movieId" binding:"required"`
	Date        string   `json:"date" binding:"required"`
	Time        string   `json:"time" binding:"required"`
	Hall        string   `json:"hall" binding:"required"`
	ReservSeats []string `json:"reservSeats" binding:"required"`
}

func (server *Server) AddReservation(ctx *gin.Context) {
	var req reservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dateValue, err := util.ParseDate(req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing date"})
		return
	}
	movie, err := server.store.Movie.GetMovie(ctx, req.MovieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := server.store.Users.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var repertoire models.Repertoire
	repertoire, err = server.store.Repertoire.GetRepertoireByMovieDateTimeHall(ctx, req.MovieID, dateValue, req.Time, req.Hall)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	numOfSeats := len(req.ReservSeats)

	if repertoire.NumOfTickets < repertoire.NumOfResTickets+numOfSeats {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not repertoire by movieId, date, time and hall"})
		return

	}

	updatedSeats := append(repertoire.ReservSeats, req.ReservSeats...)

	update := bson.M{
		"$set": bson.M{
			"numOfResTickets": repertoire.NumOfResTickets + numOfSeats,
			"reservSeats":     updatedSeats,
		},
	}

	opts := options.Update().SetUpsert(false)
	_, err = server.store.Repertoire.Col.UpdateOne(ctx, bson.M{"_id": repertoire.ID}, update, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reservation := models.Reservation{
		Username:      user.Username,
		UserID:        user.ID,
		MovieID:       movie.ID,
		MovieTitle:    movie.Title,
		RepertoiresID: repertoire.ID,
		Date:          repertoire.Date,
		Time:          repertoire.Time,
		Hall:          repertoire.Hall,
		CreationDate:  time.Now(),
		ReservSeats:   req.ReservSeats,
	}

	_, err = server.store.Reservation.Col.InsertOne(ctx, reservation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reservation added successfully"})
}
