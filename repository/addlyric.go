package repository

import (
	"context"
	"log"
	"lyrical-app/database"
	"lyrical-app/graph/model"
	"lyrical-app/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertLyric(d *database.DB, content, songID *string) (*model.Song, error) {
	collection := d.Client.Database(SongDB).Collection(LyricsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lyric := Lyric{}
	lyric.ID = primitive.NewObjectID()
	lyric.Content = *content
	lyric.Likes = 0
	lyric.SongID = utils.ToObjectID(*songID)
	_, err := collection.InsertOne(ctx, lyric)
	if err != nil {
		log.Println("------Error on InsertOne:", err)
		return nil, err
	}

	// id := fmt.Sprintf("%s", res.InsertedID.(primitive.ObjectID).Hex())

	// log.Println("++++++Lyric ID:", id)
	// return &model.Lyric{
	// 	ID:      &id,
	// 	Content: content,
	// 	SongID:  songID,
	// }, nil
	// msong := model.Song{}
	// msongId := "12345"
	// msongTitle := "this is a test"
	// msong.ID = &msongId
	// msong.Title = &msongTitle
	return GetSong(d, *songID)
}
