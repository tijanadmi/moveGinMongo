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
	Col *mongo.Collection
}

func (c *MovieClient) InitMovies(ctx context.Context) {
	setupIndexes(ctx, c.Col, "title")
}

// AddMovie adds a new movie to the MongoDB collection
func (c *MovieClient) AddMovie(ctx context.Context, movie *models.Movie) error {
	movie.ID = primitive.NewObjectID()
	_, err := c.Col.InsertOne(ctx, movie)
	if err != nil {
		log.Print(fmt.Errorf("could not add new movie: %w", err))
		return err
	}
	return nil
}

// ListMovies returns all movies from the MongoDB collection
func (c *MovieClient) ListMovies(ctx context.Context) ([]models.Movie, error) {
	movies := make([]models.Movie, 0)
	cur, err := c.Col.Find(ctx, bson.M{})
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
	res := c.Col.FindOne(ctx, bson.M{"_id": objID})
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
	res, err := c.Col.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{
		{"$set", bson.D{
			{"title", movie.Title},
			{"duration", movie.Duration},
			{"genre", movie.Genre},
			{"directors", movie.Directors},
			{"actors", movie.Actors},
			{"screening", movie.Screening},
			{"plot", movie.Plot},
			{"poster", movie.Poster},
			{"screenings", movie.Screenings},
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
	res, err := c.Col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the movie with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.DeletedCount), nil
}

// GetHall returns a hall by ID from the MongoDB collection
func (c *MovieClient) SearchMovies(ctx context.Context, movieId string) ([]models.Movie, error) {
	movies := make([]models.Movie, 0)
	fmt.Println(movieId)
	// Provera inicijalizacije kolekcije
	if c.Col == nil {
		log.Print(fmt.Errorf("collection is not initialized:"))
		return nil, fmt.Errorf("collection is not initialized")
	}

	// Dinamičko kreiranje match stage-a

	var matchStage bson.D
	if movieId != "0" {
		objectId, err := primitive.ObjectIDFromHex(movieId)
		if err != nil {
			log.Print(fmt.Errorf("invalid movie ID: %w", err))
			return nil, err
		}
		matchStage = bson.D{{"$match", bson.D{{"_id", objectId}}}}
		//matchStage = bson.D{{"$match", bson.D{{"_id", movieId}}}}
	} else {
		matchStage = bson.D{{"$match", bson.D{}}}
	}

	pipeline := mongo.Pipeline{
		matchStage,
		{
			{"$lookup", bson.D{
				{"from", "repertoires"},
				{"localField", "_id"},
				{"foreignField", "movieId"},
				{"as", "screenings"},
			}},
		},
		{
			{"$project", bson.D{
				{"_id", 1},
				{"title", 1},
				{"duration", 1},
				{"genre", 1},
				{"directors", 1},
				{"actors", 1},
				{"screening", 1},
				{"plot", 1},
				{"poster", 1},
				{"screenings.date", 1},
				{"screenings.time", 1},
				{"screenings.hall", 1},
			}},
		},
	}

	// Izvršavanje agregacije
	cursor, err := c.Col.Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(fmt.Errorf("could not aggregate movies: %w", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	// Parsiranje rezultata
	if err := cursor.All(ctx, &movies); err != nil {
		log.Print(fmt.Errorf("could not unmarshal the movies results: %w", err))
		return nil, err
	}

	return movies, nil

}
