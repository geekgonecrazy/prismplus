package store

import (
	"errors"

	"github.com/geekgonecrazy/prismplus/models"
)

//Store is an interface that the storage implementers should implement
type Store interface {
	CreateStreamer(streamer *models.Streamer) error
	GetStreamers() ([]models.Streamer, error)
	GetStreamerByID(id int) (models.Streamer, error)
	GetStreamerByStreamKey(key string) (models.Streamer, error)
	UpdateStreamer(streamer *models.Streamer) error
	DeleteStreamer(id int) error

	CheckDb() error
}

var ErrNotFound = errors.New("record not found")
