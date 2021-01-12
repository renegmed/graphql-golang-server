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

// working
func GetSong(d *database.DB, songId string) (*model.Song, error) {
	objID, err := primitive.ObjectIDFromHex(songId)
	if err != nil {
		return &model.Song{}, err
	}

	collection := d.Client.Database(SongDB).Collection(SongsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"_id", objID},
			}},
		},
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
		return &model.Song{}, err
	}

	var songs []bson.M
	//var songs []model.Song
	if err = cursor.All(ctx, &songs); err != nil {
		log.Println("Error on cursor:", err)
		return &model.Song{}, err
	}

	//log.Println("+++ songs:\n\t", songs)
	/*
		   [map[_id:ObjectID("5ffcd62f4df29b5b3958afeb")
				lyrics:[map[_id:ObjectID("5ffce332b60460af9c4a1836")
							content:It's nice to see you again
							likes:5
							songid:ObjectID("5ffcd62f4df29b5b3958afeb")]]
				title:Hello Dolly]]

	*/
	//fmt.Println(songs[0]["title"])

	// log.Println("++++aggregate song:\n\t", song)

	msong := model.Song{}

	id := songs[0]["_id"].(primitive.ObjectID).Hex()
	title := songs[0]["title"].(string)
	msong.ID = &id
	msong.Title = &title
	lyrics := songs[0]["lyrics"].(primitive.A)
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

	return &msong, nil
}
