package store

import "time"

type Variation struct {
	VariationID string    `json:"variationId" bson:"variationId"`
	Name        string    `json:"name" bson:"name"`
	Location    string    `json:"location" bson:"location"`
	Width       int       `json:"width" bson:"width"`
	Height      int       `json:"height" bson:"height"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}
type Media struct {
	MediaID    string      `json:"mediaId" bson:"mediaId"`
	Variations []Variation `json:"variations" bson:"variations"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt" bson:"updatedAt"`
}

type CreateMediaParams struct {
	MediaID string `bson:"mediaId"`
}

type MediaStore interface {
	CreateMedia(params CreateMediaParams) (media *Media, err error)
}
