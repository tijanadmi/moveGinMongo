package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/tijanadmi/moveginmongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// HallClient is the client responsible for querying MongoDB
type HallClient struct {
	Col *mongo.Collection
}

func (c *HallClient) InitHalls(ctx context.Context) {
	setupIndexes(ctx, c.Col, "name")
}

// AddHall adds a new hall to the MongoDB collection
func (c *HallClient) InsertHall(ctx context.Context, hall *models.Hall) error {
	hall.ID = primitive.NewObjectID()
	_, err := c.Col.InsertOne(ctx, hall)
	if err != nil {
		log.Print(fmt.Errorf("could not add new hall: %w", err))
		return err
	}
	return nil
}

// ListHalls returns all halls from the MongoDB collection
func (c *HallClient) ListHalls(ctx context.Context) ([]models.Hall, error) {
	halls := make([]models.Hall, 0)
	cur, err := c.Col.Find(ctx, bson.M{})
	if err != nil {
		log.Print(fmt.Errorf("could not get all halls: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &halls); err != nil {
		log.Print(fmt.Errorf("could marshall the halls results: %w", err))
		return nil, err
	}

	return halls, nil
}

// GetHall returns a hall by ID from the MongoDB collection
func (c *HallClient) SearchHall(ctx context.Context, name string) ([]models.Hall, error) {
	halls := make([]models.Hall, 0)

	// Provera inicijalizacije kolekcije
	if c.Col == nil {
		log.Print(fmt.Errorf("collection is not initialized:"))
		return nil, fmt.Errorf("collection is not initialized")
	}

	cur, err := c.Col.Find(ctx, bson.M{"name": name})

	if err != nil {
		log.Print(fmt.Errorf("could not get all halls: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &halls); err != nil {
		log.Print(fmt.Errorf("could marshall the halls results: %w", err))
		return nil, err
	}

	return halls, nil

}

// UpdateHall updates a hall by ID in the MongoDB collection
func (c *HallClient) UpdateHall(ctx context.Context, id string, hall models.Hall) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.Col.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{
		{"$set", bson.D{
			{"name", hall.Name},
			{"rows", hall.Rows},
			{"cols", hall.Cols},
		}},
	})
	if err != nil {
		log.Print(fmt.Errorf("could not update hall with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

// DeleteHall deletes a hall by ID from the MongoDB collection
func (c *HallClient) DeleteHall(ctx context.Context, id string) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.Col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the hall with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.DeletedCount), nil
}
