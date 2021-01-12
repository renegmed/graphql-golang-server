package graph

import (
	"context"
	"lyrical-app/database"
	"lyrical-app/graph/generated"
	"lyrical-app/utils"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
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

func addSong(t *testing.T, c *client.Client, title string) {
	tpl := `mutation {
		addSong(title: "{{ .Title }}" ){
		    id 
		    title
		}	
	}
	`
	var resp struct {
		AddSong struct {
			ID    string
			Title string
		}
	}
	vars := make(map[string]interface{})
	vars["Title"] = title

	query, err := utils.CreateTemplate(tpl, vars)
	if err != nil {
		t.Fail()
	}

	c.MustPost(query, &resp)

	//log.Println("Add Song:", resp.AddSong)

	require.NotEqual(t, "", resp.AddSong.ID)
	require.Equal(t, sampleSongTitle, resp.AddSong.Title)

	lastSongInsertID = resp.AddSong.ID
}

func getAllSongs(t *testing.T, c *client.Client) {
	query := `{
		songs{
		    id 
		    title
		}	
	}
	`
	var resp struct {
		Songs []struct {
			ID    string
			Title string
		}
	}

	c.MustPost(query, &resp)

	//log.Println("Song List:", resp.Songs)

	require.NotEqual(t, 0, len(resp.Songs))

}

func getSong(t *testing.T, c *client.Client, songId string) {
	tpl := `query {
		song(id: "{{ .ID }}" ){
		    id 
			title
			lyrics {
				content
				likes
			}
		}	
	}
	`
	var resp struct {
		Song struct {
			ID     string
			Title  string
			Lyrics []struct {
				Likes   int
				Content string
			}
		}
	}
	vars := make(map[string]interface{})
	vars["ID"] = songId

	query, err := utils.CreateTemplate(tpl, vars)
	if err != nil {
		t.Fail()
	}

	c.MustPost(query, &resp)

	require.NotEqual(t, "", resp.Song.ID)
	require.Equal(t, songId, resp.Song.ID)
	require.Equal(t, "My Way", resp.Song.Title)
	require.Len(t, resp.Song.Lyrics, 4)
}

func getLyric(t *testing.T, c *client.Client, lyricID string) {
	tpl := `query {
		lyric(id: "{{ .ID }}" ){
		    id 
			likes
			content
			song {
				id 
				title
			}
		}	
	}
	`
	var resp struct {
		Lyric struct {
			ID      string
			Content string
			Likes   int
			Song    struct {
				ID    string
				Title string
			}
		}
	}
	vars := make(map[string]interface{})

	vars["ID"] = lyricID

	query, err := utils.CreateTemplate(tpl, vars)
	if err != nil {
		t.Fail()
	}

	//log.Println("++++ query lyric:", query)

	c.MustPost(query, &resp)

	//log.Println("query Lyric:", resp.Lyric)

	require.Equal(t, lyricID, resp.Lyric.ID)
	require.NotEqual(t, nil, resp.Lyric.Song)

}

func addLyricToSong(t *testing.T, c *client.Client, songID, content string) {
	tpl := `mutation {
		addLyricToSong(content: "{{ .Content }}", songId: "{{ .SongID }}" ){
		    id 
			title
			lyrics {
				likes 
				content
				song {
					id 
					title
				}
			}
		}	
	}
	`
	var resp struct {
		AddLyricToSong struct {
			ID     string
			Title  string
			Lyrics []struct {
				ID      string
				Content string
				Likes   int
				Song    struct {
					ID    string
					Title string
				}
			}
		}
	}
	vars := make(map[string]interface{})
	vars["Content"] = content
	vars["SongID"] = songID

	query, err := utils.CreateTemplate(tpl, vars)
	if err != nil {
		t.Fail()
	}

	//log.Println("++++ insert lyric query:", query)

	c.MustPost(query, &resp)

	//log.Println("Insert Lyric:", resp.AddLyricToSong)

	require.NotEqual(t, "", resp.AddLyricToSong.ID)
	require.Equal(t, sampleSongTitle, resp.AddLyricToSong.Title)

	lastLyricInsertID = resp.AddLyricToSong.Lyrics[0].ID
}

func likeLyric(t *testing.T, c *client.Client, lyricID string) {
	tpl := `mutation {
		likeLyric(id: "{{ .ID }}" ){
		    id 
			likes
			content
			song {
				id 
				title
			}
		}	
	}
	`
	var resp struct {
		LikeLyric struct {
			ID      string
			Content string
			Likes   int
			Song    struct {
				ID    string
				Title string
			}
		}
	}
	vars := make(map[string]interface{})
	vars["ID"] = lyricID

	query, err := utils.CreateTemplate(tpl, vars)
	if err != nil {
		t.Fail()
	}

	//log.Println("++++ likeLyric query:\n", query)

	c.MustPost(query, &resp)

	//log.Println("query Lyric:", resp.LikeLyric)

	require.Equal(t, lyricID, resp.LikeLyric.ID)
	require.Equal(t, 6, resp.LikeLyric.Likes)
}

func deleteSong(t *testing.T, c *client.Client, songID string) {
	tpl := `mutation {
		deleteSong(id: "{{ .ID }}" ){
		    id 
			title 
			lyrics {
				id 
				content 
				likes
			}
		}	
	}
	`
	var resp struct {
		DeleteSong struct {
			ID     string
			Title  string
			Lyrics []struct {
				ID      string
				Content string
				Likes   int
			}
		}
	}
	vars := make(map[string]interface{})
	vars["ID"] = songID

	query, err := utils.CreateTemplate(tpl, vars)
	require.NoError(t, err)

	c.MustPost(query, &resp)

	//log.Println("Add Song:", resp.DeleteSong)

	require.Equal(t, songID, resp.DeleteSong.ID)

	lastSongInsertID = ""
}
func TestSong(t *testing.T) {

	db := database.NewDb()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	utils.PrepareSongsData(ctx, t, db.Client)
	utils.PrepareLyricsData(ctx, t, db.Client)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewResolver(db)})))

	t.Run("Get all songs", func(t *testing.T) {
		getAllSongs(t, c)
	})

	t.Run("Get one song", func(t *testing.T) {
		getSong(t, c, "5ffcd62f4df29b5b3958afea") //lastInsertID)
	})

	t.Run("Get one lyric", func(t *testing.T) {
		getLyric(t, c, "5ffd215edafa9ca71a9dcd8f")
	})

	t.Run("Like a lyric", func(t *testing.T) {
		likeLyric(t, c, "5ffd215edafa9ca71a9dcd8f")
	})

	t.Run("Add song", func(t *testing.T) {
		addSong(t, c, sampleSongTitle)
	})

	t.Run("Add lyric to a song", func(t *testing.T) {
		addLyricToSong(t, c, lastSongInsertID, sampleSongContent)
	})

	t.Run("Delete added song", func(t *testing.T) {
		deleteSong(t, c, lastSongInsertID)
	})

}
