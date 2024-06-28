package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
	"github.com/tijanadmi/moveginmongo/util"
)

// GetAllReservationsForUser godoc
// @Security bearerAuth
// @Summary Get all the existing reservations for user
// @Description Get all the existing reservations for user
// @ID GetAllReservationsForUser
// @Accept  json
// @Produce  json
// @Param  username query string true "Username"
// @Success 200 {array} models.Reservations
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /reservationforuser [get]
func (server *Server) GetAllReservationsForUser(ctx *gin.Context) {

	username := ctx.Query("username")

	repertoires, err := server.store.Reservation.GetAllReservationsForUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repertoires)
}

/*type reservationRequest struct {
	Username    string   `json:"username" binding:"required"`
	MovieID     string   `json:"movieId" binding:"required"`
	Date        string   `json:"date" binding:"required"`
	Time        string   `json:"time" binding:"required"`
	Hall        string   `json:"hall" binding:"required"`
	ReservSeats []string `json:"reservSeats" binding:"required"`
}*/

// AddReservation godoc
// @Security bearerAuth
// @Summary Insert new reservation
// @Description Insert new reservation
// @ID AddReservation
// @Accept  json
// @Produce  json
// @Param reservation body models.Reservations true "Create reservation"
// @Success 201 {array} models.Reservations
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /reservation [post]
func (server *Server) AddReservation(ctx *gin.Context) {
	var req reservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: "Invalid input"})
		return
	}

	dateValue, err := util.ParseDate(req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: "Error parsing date"})
		return
	}
	movie, err := server.store.Movie.GetMovie(ctx, req.MovieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	user, err := server.store.Users.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	var repertoire models.Repertoire
	repertoire, err = server.store.Repertoire.GetRepertoireByMovieDateTimeHall(ctx, req.MovieID, dateValue, req.Time, req.Hall)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	numOfSeats := len(req.ReservSeats)

	if repertoire.NumOfTickets < repertoire.NumOfResTickets+numOfSeats {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not repertoire by movieId, date, time and hall"})
		return

	}

	updatedSeats := append(repertoire.ReservSeats, req.ReservSeats...)
	sort.Strings(updatedSeats)
	repertoire.ReservSeats = updatedSeats
	repertoire.NumOfResTickets = repertoire.NumOfResTickets + numOfSeats

	/***** Update Repertoire ***/
	modifiedCount, err := server.store.Repertoire.UpdateRepertoire(ctx, repertoire.ID.Hex(), repertoire)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, apiErrorResponse{Error: "repertoire not found"})
		return
	}

	/***** Isert Reservation ****/
	sort.Strings(req.ReservSeats)
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

	if err := server.store.Reservation.InsertReservation(ctx, &reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: "Reservation added successfully"})
}

// CancelReservation godoc
// @Security bearerAuth
// @Summary Delete a single reservation
// @Description Delete a single reservation
// @ID CancelReservation
// @Accept  json
// @Produce  json
// @Param  id path string true "reservation ID"
// @Success 200 {array} apiResponse
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /reservation/{id} [delete]
func (server *Server) CancelReservation(ctx *gin.Context) {
	id := ctx.Param("id")
	reservation, err := server.store.Reservation.GetReservationById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}
	repertoire, err := server.store.Repertoire.GetRepertoire(ctx, reservation.RepertoiresID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	difSeats := difference(repertoire.ReservSeats, reservation.ReservSeats)
	repertoire.ReservSeats = difSeats
	repertoire.NumOfResTickets = len(difSeats)

	/*** Update repertoire ****/
	modifiedCount, err := server.store.Repertoire.UpdateRepertoire(ctx, reservation.RepertoiresID.Hex(), repertoire)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}
	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, apiErrorResponse{Error: "repertoire not found for update"})
		return
	}

	/**** Delete reservation ****/
	deletedCount, err := server.store.Reservation.DeleteReservation(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, apiErrorResponse{Error: "reservation not found"})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: "Reservation canceled successfully"})

}

func difference(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	for _, item := range slice2 {
		m[item] = true
	}
	var diff []string
	for _, item := range slice1 {
		if !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
