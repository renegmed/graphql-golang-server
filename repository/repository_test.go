package repository

import (
	"context"
	"lyrical-app/database"
	"lyrical-app/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sampleSongTitle   = "One Way or Another"
	sampleSongContent = "I gonna get you"
)

var (
	lastSongInsertID  string
	lastLyricInsertID string
)

func TestRepository(t *testing.T) {
	db := database.NewDb()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	utils.PrepareSongsData(ctx, t, db.Client)
	utils.PrepareLyricsData(ctx, t, db.Client)

	t.Run("Get song with lyrics", func(t *testing.T) {
		id := "5ffcd62f4df29b5b3958afea"
		title := "My Way"
		song, err := GetSong(db, id) //"5ffcd62f4df29b5b3958afeb")
		require.NoError(t, err)
		// fmt.Println("-----------")
		// fmt.Println(*song.Title)
		// for i, _ := range song.Lyrics {
		// 	fmt.Println(*(*song.Lyrics[i]).Content, "  ", *(*song.Lyrics[i]).Likes)
		// 	fmt.Println("Song title:", *(*song.Lyrics[i]).Song.Title)
		// 	fmt.Println(".")
		// }
		// fmt.Println("-----------")

		require.Equal(t, *song.ID, id)
		require.Equal(t, *song.Title, title)
		require.Len(t, song.Lyrics, 4)

	})

	t.Run("Get all songs with lyrics", func(t *testing.T) {
		// id := "5ffcd62f4df29b5b3958afea"
		// title := "My Way"
		songs, err := GetSongs(db) //"5ffcd62f4df29b5b3958afeb")
		require.NoError(t, err)

		// fmt.Println("-----------")
		// fmt.Println(*song.Title)
		// for i, _ := range song.Lyrics {
		// 	fmt.Println(*(*song.Lyrics[i]).Content, "  ", *(*song.Lyrics[i]).Likes)
		// 	fmt.Println("Song title:", *(*song.Lyrics[i]).Song.Title)
		// 	fmt.Println(".")
		// }
		// fmt.Println("-----------")
		require.Len(t, songs, 2)

	})

	t.Run("Get a lyric", func(t *testing.T) {
		id := "5ffd215edafa9ca71a9dcd8f"
		likes := 5
		songID := "5ffcd62f4df29b5b3958afeb"
		title := "Hello Dolly"

		lyric, err := GetLyric(db, id) //"5ffcd62f4df29b5b3958afeb")
		require.NoError(t, err)
		require.Equal(t, id, *lyric.ID)
		require.Equal(t, likes, *lyric.Likes)
		require.Equal(t, songID, *(*lyric.Song).ID)
		require.Equal(t, title, *(*lyric.Song).Title)
	})

	t.Run("Like a lyric", func(t *testing.T) {
		id := "5ffd215edafa9ca71a9dcd8f"
		songID := "5ffcd62f4df29b5b3958afeb"
		songTitle := "Hello Dolly"

		lyric, err := GetLyric(db, id)
		require.NoError(t, err)

		currentLikes := *lyric.Likes

		newLyric, err := LikeLyric(db, lyric) //"5ffcd62f4df29b5b3958afeb")
		require.NoError(t, err)

		require.Equal(t, currentLikes+1, *newLyric.Likes)
		require.Equal(t, songID, *newLyric.Song.ID)
		require.Equal(t, songTitle, *newLyric.Song.Title)
	})

	t.Run("Add a new song", func(t *testing.T) {
		title := sampleSongTitle
		song, err := SaveSong(db, &title)
		require.NoError(t, err)

		require.NotEqual(t, nil, *song.ID)
		require.Equal(t, sampleSongTitle, *song.Title)

		lastSongInsertID = *song.ID
	})

	t.Run("Add lyric to new song", func(t *testing.T) {
		songID := lastSongInsertID
		content := sampleSongContent
		song, err := InsertLyric(db, &content, &songID)
		require.NoError(t, err)

		require.Equal(t, songID, *song.ID)
		require.Equal(t, sampleSongTitle, *song.Title)

		require.Equal(t, content, *(*song).Lyrics[0].Content)
	})

	t.Run("Get all songs after addition", func(t *testing.T) {
		songs, err := GetSongs(db)
		require.NoError(t, err)
		require.Len(t, songs, 3)
	})

	t.Run("Delete the new song", func(t *testing.T) {

		songDeleted, err := DeleteSong(db, &lastSongInsertID)
		require.NoError(t, err)

		require.Equal(t, lastSongInsertID, *songDeleted.ID)

	})
}
