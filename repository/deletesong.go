package repository

import (
	"context"
	"errors"
	"fmt"
	"lyrical-app/database"
	"lyrical-app/graph/model"
	"lyrical-app/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteSong(d *database.DB, songID *string) (*model.Song, error) {
	collection := d.Client.Database(SongDB).Collection(SongsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	song, err := GetSong(d, *songID)
	if err != nil {
		return &model.Song{}, err
	}

	err = deleteSongLyrics(d, songID)
	if err != nil {
		return &model.Song{}, err
	}
	result, err := collection.DeleteOne(
		ctx,
		bson.D{
			{"_id", utils.ToObjectID(*songID)},
		},
	)
	if result.DeletedCount == 0 {
		return &model.Song{}, errors.New(fmt.Sprintf("Song %s has not been deleted", *songID))
	}
	return song, nil
}

func deleteSongLyrics(d *database.DB, songID *string) error {
	collection := d.Client.Database(SongDB).Collection(LyricsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteMany(
		ctx,
		bson.D{
			{"songid", utils.ToObjectID(*songID)},
		},
	)
	if err != nil {
		return err
	}
	// if result.DeletedCount == 0 {
	// 	return errors.New(fmt.Sprintf("Lyrics for song %s has not been deleted", *songID))
	// }
	return nil
}
