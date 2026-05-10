package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoRepo struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoRepo(ctx context.Context, uri, dbName string) (*MongoRepo, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctxPing, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := client.Ping(ctxPing, readpref.Primary()); err != nil {
		return nil, err
	}

	coll := client.Database(dbName).Collection("files")
	return &MongoRepo{
		client: client,
		coll:   coll,
	}, nil
}

func (m *MongoRepo) InsertRecord(ctx context.Context, record *FileRecord) error {
	_, err := m.coll.InsertOne(ctx, record)
	return err
}

func (m *MongoRepo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
