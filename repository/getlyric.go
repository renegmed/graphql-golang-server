package repository

import (
	"context"
	"log"
	"lyrical-app/database"
	"lyrical-app/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetLyric(d *database.DB, lyricId string) (*model.Lyric, error) {

	objID, err := primitive.ObjectIDFromHex(lyricId)
	if err != nil {
		return nil, err
	}
	collection := d.Client.Database(SongDB).Collection(LyricsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"_id", objID},
			}},
		},
		{
			{"$lookup", bson.D{
				{"from", "songs"},
				{"localField", "songid"},
				{"foreignField", "_id"},
				{"as", "songs"},
			}},
		},
		{
			{"$project", bson.D{
				{"songid", 0},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error on aggregate:", err)
		log.Println(err)
		return &model.Lyric{}, err
	}

	var lyrics []bson.M
	if err = cursor.All(ctx, &lyrics); err != nil {
		log.Println("Error on cursor:", err)
		return &model.Lyric{}, err
	}

	// log.Println("+++++ lyric:", lyrics)
	/*
		[map[_id:ObjectID("5ffd215edafa9ca71a9dcd8f")
			content:It's nice to see you again
			likes:5
			songs:[map[_id:ObjectID("5ffcd62f4df29b5b3958afeb")
			           title:Hello Dolly]]]]
	*/
	mLyric := model.Lyric{}
	id := lyrics[0]["_id"].(primitive.ObjectID).Hex()
	content := lyrics[0]["content"].(string)
	likes := int((lyrics[0]["likes"].(int32)))

	mLyric.ID = &id
	mLyric.Content = &content
	mLyric.Likes = &likes

	songs := lyrics[0]["songs"].(primitive.A)

	for i, _ := range songs {
		songID := songs[i].(primitive.M)["_id"].(primitive.ObjectID).Hex()
		songTitle := songs[i].(primitive.M)["title"].(string)

		song := model.Song{}
		song.ID = &songID
		song.Title = &songTitle

		mLyric.Song = &song
		break
	}
	return &mLyric, nil
}
