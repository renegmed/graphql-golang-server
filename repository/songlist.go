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

func GetSongs(d *database.DB) ([]*model.Song, error) {

	collection := d.Client.Database(SongDB).Collection(SongsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{"$lookup", bson.D{
				{"from", "lyrics"},
				{"localField", "_id"},
				{"foreignField", "songid"},
				{"as", "lyrics"},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error on aggregate:", err)
		log.Println(err)
		return nil, err
	}

	var songs []bson.M
	if err = cursor.All(ctx, &songs); err != nil {
		log.Println("Error on cursor:", err)
		return nil, err
	}

	// log.Println("+++ songs:\n\t", songs)
	/*
		   [map[_id:ObjectID("5ffcd62f4df29b5b3958afeb")
				lyrics:[map[_id:ObjectID("5ffce332b60460af9c4a1836")
							content:It's nice to see you again
							likes:5
							songid:ObjectID("5ffcd62f4df29b5b3958afeb")]]
				title:Hello Dolly]]

	*/

	msongs := []*model.Song{}

	for k, _ := range songs {
		id := songs[k]["_id"].(primitive.ObjectID).Hex()
		title := songs[k]["title"].(string)
		msong := model.Song{}
		msong.ID = &id
		msong.Title = &title
		lyrics := songs[k]["lyrics"].(primitive.A)
		for i, _ := range lyrics {
			//fmt.Println(lyrics[i].(primitive.M)["content"])
			lyric := model.Lyric{}
			var lyricId = lyrics[i].(primitive.M)["_id"].(primitive.ObjectID).Hex()
			var likes = int((lyrics[i].(primitive.M)["likes"].(int32)))
			var content = lyrics[i].(primitive.M)["content"].(string)
			lyric.ID = &lyricId
			lyric.Likes = &likes
			lyric.Content = &content
			lyric.Song = &msong
			msong.Lyrics = append(msong.Lyrics, &lyric)
		}
		msongs = append(msongs, &msong)
	}
	return msongs, nil
}
