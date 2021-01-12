package repository

import (
	"context"
	"fmt"
	"lyrical-app/database"
	"lyrical-app/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveSong(d *database.DB, title *string) (*model.Song, error) {
	collection := d.Client.Database(SongDB).Collection(SongsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Song{}
	s.ID = primitive.NewObjectID()
	s.Title = *title
	res, err := collection.InsertOne(ctx, s)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s", res.InsertedID.(primitive.ObjectID).Hex())
	return &model.Song{
		ID:    &id,
		Title: title,
	}, nil
}
