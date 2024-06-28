package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tijanadmi/moveginmongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RepertoireClient is the client responsible for querying mongodb
type ReservationClient struct {
	Col *mongo.Collection
}

// AddReservation adds a new reservation to the MongoDB collection
func (c *ReservationClient) InsertReservation(ctx context.Context, reservation *models.Reservation) error {
	reservation.ID = primitive.NewObjectID()

	_, err := c.Col.InsertOne(ctx, reservation)
	if err != nil {
		log.Print(fmt.Errorf("could not add new reservation: %w", err))
		return err
	}
	return nil
}

// GetReservationById returns a reservations based on its ID
func (c *ReservationClient) GetReservationById(ctx context.Context, id string) (models.Reservation, error) {
	var reservation models.Reservation
	objID, _ := primitive.ObjectIDFromHex(id)
	res := c.Col.FindOne(ctx, bson.M{"_id": objID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return reservation, nil
		}
		log.Print(fmt.Errorf("error when finding the repertoire [%s]: %q", id, res.Err()))
		return reservation, res.Err()
	}

	if err := res.Decode(&reservation); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", id, err))
		return reservation, err
	}
	return reservation, nil
}

// GetReservation returns a all reservation based on username
func (c *ReservationClient) GetAllReservationsForUser(ctx context.Context, username string) ([]models.Reservation, error) {
	reservations := make([]models.Reservation, 0)

	cur, err := c.Col.Find(ctx, bson.M{"username": username})
	if err != nil {
		log.Print(fmt.Errorf("could not get all reservations [%s]: %w", username, err))
		return nil, err
	}

	if err = cur.All(ctx, &reservations); err != nil {
		log.Print(fmt.Errorf("could marshall the repertoires results: %w", err))
		return nil, err
	}

	return reservations, nil
}

// DeleteReservation deletes a reservation based on its ID
func (c *ReservationClient) DeleteReservation(ctx context.Context, id string) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.Col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the repertoire with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.DeletedCount), nil
}
