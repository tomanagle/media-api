package store

import (
	"context"
	"time"
)

type Variation struct {
	Name     string `json:"name" bson:"name"`
	Location string `json:"location" bson:"location"`
	Width    int    `json:"width" bson:"width"`
	Height   int    `json:"height" bson:"height"`
}
type Media struct {
	MediaID    string      `json:"mediaId" bson:"mediaId"`
	Variations []Variation `json:"variations" bson:"variations"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt" bson:"updatedAt"`
}

type CreateMediaParams struct {
	MediaID    string      `bson:"mediaId"`
	Variations []Variation `bson:"variations"`
}

type MediaStore interface {
	CreateMedia(ctx context.Context, params CreateMediaParams) (media *Media, err error)
	GetMedia(ctx context.Context, params GetMediaParams) (media []Media, err error)
}

type GetMediaParams struct {
	Skip  int64 `bson:"skip"`
	Limit int64 `bson:"limit"`
}
