package appinit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetMongoClientSuccess(t *testing.T) {
	expected := &mongo.Client{}
	old := mongoConnect
	mongoConnect = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
		return expected, nil
	}
	defer func() { mongoConnect = old }()

	client := GetMongoClient(options.Client())
	assert.Equal(t, expected, client)
}

func TestGetMongoClientError(t *testing.T) {
	old := mongoConnect
	mongoConnect = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
		return nil, errors.New("fail")
	}
	defer func() { mongoConnect = old }()

	assert.PanicsWithValue(t, "Mongo Connection Error!", func() {
		GetMongoClient(options.Client())
	})
}
