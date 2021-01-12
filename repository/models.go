package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

const SongDB = "songdb"
const SongsCollection = "songs"
const LyricsCollection = "lyrics"

type Song struct {
	ID    primitive.ObjectID `bson:"_id"`
	Title string             `bson:"title"`
}

type Lyric struct {
	ID      primitive.ObjectID `bson:"_id"`
	Likes   int                `bson:"likes"`
	Content string             `bson:"content"`
	SongID  primitive.ObjectID `bson:"songid"`
}
