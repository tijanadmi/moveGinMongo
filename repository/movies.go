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

// MovieModel sa CRUD operacijama
// MovieClient is the client responsible for querying MongoDB
type MovieClient struct {
    col  *mongo.Collection
}


func (c *MovieClient) InitMovies(ctx context.Context) {
    setupIndexes(ctx, c.col, "title")
}

// AddMovie adds a new movie to the MongoDB collection
func (c *MovieClient) AddMovie(ctx context.Context, movie *models.Movie) error {
    movie.ID = primitive.NewObjectID()
    _, err := c.col.InsertOne(ctx, movie)
    if err != nil {
        log.Print(fmt.Errorf("could not add new movie: %w", err))
        return err
    }
    return nil
}

// ListMovies returns all movies from the MongoDB collection
func (c *MovieClient) ListMovies(ctx context.Context) ([]models.Movie, error) {
    movies := make([]models.Movie, 0)
    cur, err := c.col.Find(ctx, bson.M{})
    if err != nil {
        log.Print(fmt.Errorf("could not get all movies: %w", err))
        return nil, err
    }

    if err = cur.All(ctx, &movies); err != nil {
        log.Print(fmt.Errorf("could not marshall the movies results: %w", err))
        return nil, err
    }

    return movies, nil
}

// GetMovie returns a movie by ID from the MongoDB collection
func (c *MovieClient) GetMovie(ctx context.Context, id string) (models.Movie, error) {
    var movie models.Movie
    objID, _ := primitive.ObjectIDFromHex(id)
    res := c.col.FindOne(ctx, bson.M{"_id": objID})
    if res.Err() != nil {
        if res.Err() == mongo.ErrNoDocuments {
            return movie, nil
        }
        log.Print(fmt.Errorf("error when finding the movie [%s]: %q", id, res.Err()))
        return movie, res.Err()
    }

    if err := res.Decode(&movie); err != nil {
        log.Print(fmt.Errorf("error decoding [%s]: %q", id, err))
        return movie, err
    }
    return movie, nil
}

// UpdateMovie updates a movie by ID in the MongoDB collection
func (c *MovieClient) UpdateMovie(ctx context.Context, id string, movie models.Movie) (int, error) {
    objID, _ := primitive.ObjectIDFromHex(id)
    res, err := c.col.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{
        {"$set", bson.D{
            {"title", movie.Title},
            {"duration", movie.Duration},
            {"genre", movie.Genre},
            {"directors", movie.Directors},
            {"actors", movie.Actors},
            {"screening", movie.Screening},
            {"plot", movie.Plot},
            {"poster", movie.Poster},
            {"repertoires", movie.Repertoires},
        }},
    })
    if err != nil {
        log.Print(fmt.Errorf("could not update movie with id [%s]: %w", id, err))
        return 0, err
    }

    return int(res.ModifiedCount), nil
}

// DeleteMovie deletes a movie by ID from the MongoDB collection
func (c *MovieClient) DeleteMovie(ctx context.Context, id string) (int, error) {
    objID, _ := primitive.ObjectIDFromHex(id)
    res, err := c.col.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        log.Print(fmt.Errorf("error deleting the movie with id [%s]: %w", id, err))
        return 0, err
    }

    return int(res.DeletedCount), nil
}
