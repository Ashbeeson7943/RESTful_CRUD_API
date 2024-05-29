package data

import "errors"

type Album struct {
	ID     string  `json:"id"`
	TITLE  string  `json:"title"`
	ARTIST string  `json:"artist"`
	PRICE  float32 `json:"price"`
}

var Albums = []Album{
	{ID: "1", TITLE: "Blue Train", ARTIST: "John Coltrane", PRICE: 56.99},
	{ID: "2", TITLE: "Jeru", ARTIST: "Gerry Mulligan", PRICE: 17.99},
	{ID: "3", TITLE: "Sarah Vaughan and Clifford Brown", ARTIST: "Sarah Vaughan", PRICE: 39.99},
}

type ApiKey struct {
	OWNER string
	VALUE string
	VALID bool
}

func (k *ApiKey) IsValid() error {
	if !k.VALID {
		return errors.New("api key is no longer valid")
	} else {
		return nil
	}
}
