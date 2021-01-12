package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"lyrical-app/graph/generated"
	"lyrical-app/graph/model"
	"lyrical-app/repository"
)

func (r *mutationResolver) AddSong(ctx context.Context, title *string) (*model.Song, error) {
	return repository.SaveSong(r.Db, title)
}

func (r *mutationResolver) AddLyricToSong(ctx context.Context, content *string, songID *string) (*model.Song, error) {
	return repository.InsertLyric(r.Db, content, songID)
}

func (r *mutationResolver) LikeLyric(ctx context.Context, id *string) (*model.Lyric, error) {
	lyric, err := repository.GetLyric(r.Db, *id)
	if err != nil {
		return &model.Lyric{}, err
	}
	return repository.LikeLyric(r.Db, lyric)
}

func (r *mutationResolver) DeleteSong(ctx context.Context, id *string) (*model.Song, error) {
	return repository.DeleteSong(r.Db, id)
}

func (r *queryResolver) Songs(ctx context.Context) ([]*model.Song, error) {
	return repository.GetSongs(r.Db)
}

func (r *queryResolver) Song(ctx context.Context, id string) (*model.Song, error) {
	return repository.GetSong(r.Db, id)
}

func (r *queryResolver) Lyric(ctx context.Context, id string) (*model.Lyric, error) {
	return repository.GetLyric(r.Db, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
