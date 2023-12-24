package dbstore

import (
	"context"
	"go-media/internal/store"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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

func (s *MediaStore) CreateMedia(params store.CreateMediaParams) (media *store.Media, err error) {

	now := time.Now()

	media = &store.Media{
		MediaID:    params.MediaID,
		Variations: params.Variations,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	res, err := s.collection.InsertOne(context.Background(), media)

	if err != nil {
		return nil, err
	}

	return media, nil
}
