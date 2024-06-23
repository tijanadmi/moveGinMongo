package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Movie predstavlja podatke o filmu
type Movie struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Title       string             `bson:"title,omitempty" json:"title,omitempty"`
    Duration    string             `bson:"duration,omitempty" json:"duration,omitempty"`
    Genre       string             `bson:"genre,omitempty" json:"genre,omitempty"`
    Directors   string             `bson:"directors,omitempty" json:"directors,omitempty"`
    Actors      string             `bson:"actors,omitempty" json:"actors,omitempty"`
    Screening   time.Time          `bson:"screening,omitempty" json:"screening,omitempty"`
    Plot        string             `bson:"plot,omitempty" json:"plot,omitempty"`
    Poster      string             `bson:"poster,omitempty" json:"poster,omitempty"`
    Repertoires []Repertoire  `bson:"repertoires,omitempty" json:"repertoires,omitempty"`
}
