package dbstore

import (
	"context"
	"go-media/internal/store"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MediaStore struct {
	collection *mongo.Collection
}

type NewMediaStoreParams struct {
	Collection *mongo.Collection
}

func NewMediaStore(p NewMediaStoreParams) *MediaStore {
	if p.Collection == nil {
		panic("collection is nil")
	}
	return &MediaStore{
		collection: p.Collection,
	}
}

func (s *MediaStore) CreateMedia(ctx context.Context, params store.CreateMediaParams) (media *store.Media, err error) {

	now := time.Now()

	media = &store.Media{
		MediaID:    params.MediaID,
		Variations: params.Variations,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err = s.collection.InsertOne(ctx, media)

	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s *MediaStore) GetMedia(ctx context.Context, params store.GetMediaParams) (media []store.Media, err error) {

	opts := options.Find().SetSkip(params.Skip).SetLimit(params.Limit)

	cursor, err := s.collection.Find(ctx, bson.M{}, opts)

	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &media)

	if err != nil {
		return nil, err
	}

	return media, nil

}
