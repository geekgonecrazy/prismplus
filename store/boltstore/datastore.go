package boltstore

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/geekgonecrazy/prismplus/store"
	bolt "go.etcd.io/bbolt"
)

type boltStore struct {
	*bolt.DB
}

var (
	streamersBucket = []byte("streamers")
)

//New creates a new bolt store
func New(dataPath string) (store.Store, error) {
	db, err := bolt.Open(fmt.Sprintf("%s%s", dataPath, "data.bbolt"), 0600, &bolt.Options{Timeout: 15 * time.Second})
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists(streamersBucket); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &boltStore{db}, nil
}

func (s *boltStore) CheckDb() error {
	tx, err := s.Begin(false)
	if err != nil {
		return err
	}

	return tx.Rollback()
}

//itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
