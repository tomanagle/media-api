package nanoid

import gonanoid "github.com/matoous/go-nanoid/v2"

type NanoID struct{}

func New() *NanoID {
	return &NanoID{}
}

func (n *NanoID) New() (string, error) {
	return gonanoid.New()
}
