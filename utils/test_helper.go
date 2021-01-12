package utils

import (
	"bytes"
	"context"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ToObjectID(s string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return id
}

func ToString(objId primitive.ObjectID) string {
	id := objId.Hex()
	return id
}

func CreateTemplate(body string, vars interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(body)
	if err != nil {
		//utils.Error(err)
		return "", nil
	}
	return process(tmpl, vars)
}

func process(t *template.Template, vars interface{}) (string, error) {
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return "", err
	}
	return b.String(), nil

}

func PrepareSongsData(ctx context.Context, t *testing.T, client *mongo.Client) {
	coll := client.Database("songdb").Collection("songs")
	err := coll.Drop(context.Background())
	require.NoError(t, err)

	docs := []interface{}{
		bson.D{
			{"_id", ToObjectID("5ffcd62f4df29b5b3958afea")},
			{"title", "My Way"},
		},
		bson.D{
			{"_id", ToObjectID("5ffcd62f4df29b5b3958afeb")},
			{"title", "Hello Dolly"},
		},
	}
	result, collErr := coll.InsertMany(ctx, docs)
	require.NoError(t, collErr)
	require.Len(t, result.InsertedIDs, 2)

}

func PrepareLyricsData(ctx context.Context, t *testing.T, client *mongo.Client) {
	coll := client.Database("songdb").Collection("lyrics")
	err := coll.Drop(context.Background())
	require.NoError(t, err)

	docs := []interface{}{
		bson.D{
			{"_id", ToObjectID("5ffd215edafa9ca71a9dcd8b")},
			{"likes", 1},
			{"content", "And now"},
			{"songid", ToObjectID("5ffcd62f4df29b5b3958afea")},
		},
		bson.D{
			{"_id", ToObjectID("5ffd215edafa9ca71a9dcd8c")},
			{"likes", 2},
			{"content", "The end is near"},
			{"songid", ToObjectID("5ffcd62f4df29b5b3958afea")},
		},
		bson.D{
			{"_id", ToObjectID("5ffd215edafa9ca71a9dcd8d")},
			{"likes", 3},
			{"content", "It's time to face"},
			{"songid", ToObjectID("5ffcd62f4df29b5b3958afea")},
		},
		bson.D{
			{"_id", ToObjectID("5ffd215edafa9ca71a9dcd8e")},
			{"likes", 4},
			{"content", "The final curtain"},
			{"songid", ToObjectID("5ffcd62f4df29b5b3958afea")},
		},
		bson.D{
			{"_id", ToObjectID("5ffd215edafa9ca71a9dcd8f")},
			{"likes", 5},
			{"content", "It's nice to see you again"},
			{"songid", ToObjectID("5ffcd62f4df29b5b3958afeb")},
		},
	}

	result, collErr := coll.InsertMany(ctx, docs)
	require.NoError(t, collErr)
	require.Len(t, result.InsertedIDs, 5)
}
