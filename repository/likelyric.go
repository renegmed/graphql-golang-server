package repository

import (
	"context"
	"errors"
	"lyrical-app/database"
	"lyrical-app/graph/model"
	"lyrical-app/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func LikeLyric(d *database.DB, lyric *model.Lyric) (*model.Lyric, error) {

	// log.Println("Lyric ID:", *lyric.ID)
	// log.Println("Lyric Likes:", *lyric.Likes)

	*lyric.Likes++
	// log.Println("NEW Lyric Likes:", *lyric.Likes)

	collection := d.Client.Database(SongDB).Collection(LyricsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": utils.ToObjectID(*lyric.ID)},
		bson.D{
			{"$set", bson.D{
				{"likes", *lyric.Likes},
			}},
		},
	)
	if err != nil {
		return &model.Lyric{}, err
	}

	if result.ModifiedCount == 0 {
		return &model.Lyric{}, errors.New("Likes was not registered.")
	}

	//log.Println("+++++ udpate modified count", result.ModifiedCount)

	return lyric, nil
}
