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

// AddRepertoire adds a new repertoire to the MongoDB collection
func (c *ReservationClient) InsertReservation(ctx context.Context, reservation *models.Reservation) error {
	reservation.ID = primitive.NewObjectID()

	_, err := c.Col.InsertOne(ctx, reservation)
	if err != nil {
		log.Print(fmt.Errorf("could not add new reservation: %w", err))
		return err
	}
	return nil
}

// ListRepertoires returns all repertoires from the MongoDB collection
func (c *ReservationClient) ListRepertoires(ctx context.Context) ([]models.Repertoire, error) {
	repertoires := make([]models.Repertoire, 0)
	cur, err := c.Col.Find(ctx, bson.M{})
	if err != nil {
		log.Print(fmt.Errorf("could not get all repertoires: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &repertoires); err != nil {
		log.Print(fmt.Errorf("could marshall the repertoires results: %w", err))
		return nil, err
	}

	return repertoires, nil
}

// GetRepertoire returns a repertoire based on its ID
func (c *ReservationClient) GetRepertoire(ctx context.Context, id string) (models.Repertoire, error) {
	var repertoire models.Repertoire
	objID, _ := primitive.ObjectIDFromHex(id)
	res := c.Col.FindOne(ctx, bson.M{"_id": objID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return repertoire, nil
		}
		log.Print(fmt.Errorf("error when finding the repertoire [%s]: %q", id, res.Err()))
		return repertoire, res.Err()
	}

	if err := res.Decode(&repertoire); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", id, err))
		return repertoire, err
	}
	return repertoire, nil
}

// GetRepertoire returns a repertoires based on its movieId
func (c *ReservationClient) GetAllRepertoireForMovie(ctx context.Context, movieId string) ([]models.Repertoire, error) {
	repertoires := make([]models.Repertoire, 0)

	movieID, _ := primitive.ObjectIDFromHex(movieId)

	cur, err := c.Col.Find(ctx, bson.M{"movieId": movieID})
	if err != nil {
		log.Print(fmt.Errorf("could not get all repertoires [%s]: %w", movieId, err))
		return nil, err
	}

	if err = cur.All(ctx, &repertoires); err != nil {
		log.Print(fmt.Errorf("could marshall the repertoires results: %w", err))
		return nil, err
	}

	return repertoires, nil
}

// UpdateRepertoire updates a repertoire based on its ID
func (c *ReservationClient) UpdateRepertoire(ctx context.Context, id string, repertoire models.Repertoire) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.Col.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{
		{"$set", bson.D{
			{"movieId", repertoire.MovieID},
			{"date", repertoire.Date},
			{"time", repertoire.Time},
			{"hall", repertoire.Hall},
			{"numOfTickets", repertoire.NumOfTickets},
			{"numOfResTickets", repertoire.NumOfResTickets},
			{"reservSeats", repertoire.ReservSeats},
		}},
	})
	if err != nil {
		log.Print(fmt.Errorf("could not update repertoire with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

// DeleteRepertoire deletes a repertoire based on its ID
func (c *ReservationClient) DeleteRepertoire(ctx context.Context, id string) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.Col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the repertoire with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.DeletedCount), nil
}

// DeleteRepertoire deletes a repertoire based on its ID
func (c *ReservationClient) DeleteRepertoireForMovie(ctx context.Context, movieId string) (int, error) {
	movieID, _ := primitive.ObjectIDFromHex(movieId)
	res, err := c.Col.DeleteMany(ctx, bson.M{"movieId": movieID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the repertoire with id [%s]: %w", movieId, err))
		return 0, err
	}

	return int(res.DeletedCount), nil
}
