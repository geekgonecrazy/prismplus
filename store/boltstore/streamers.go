package boltstore

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/geekgonecrazy/prismplus/models"
	"github.com/geekgonecrazy/prismplus/store"
	bolt "go.etcd.io/bbolt"
)

func (s *boltStore) GetStreamers() ([]models.Streamer, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(streamersBucket).Cursor()

	streamers := make([]models.Streamer, 0)
	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Streamer
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		streamers = append(streamers, i)
	}

	return streamers, nil
}

func (s *boltStore) GetStreamerByID(id int) (streamer models.Streamer, err error) {
	tx, err := s.Begin(false)
	if err != nil {
		return streamer, err
	}
	defer tx.Rollback()

	bytes := tx.Bucket(streamersBucket).Get(itob(id))
	if bytes == nil {
		return streamer, store.ErrNotFound
	}

	var i models.Streamer
	if err := json.Unmarshal(bytes, &i); err != nil {
		return streamer, err
	}

	return i, nil
}

func (s *boltStore) GetStreamerByStreamKey(key string) (streamer models.Streamer, err error) {
	tx, err := s.Begin(false)
	if err != nil {
		return streamer, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(streamersBucket).Cursor()

	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Streamer
		if err := json.Unmarshal(data, &i); err != nil {
			return streamer, err
		}

		if i.StreamKey == key {
			return i, nil
		}
	}

	return streamer, store.ErrNotFound
}

func (s *boltStore) CreateStreamer(streamer *models.Streamer) error {
	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(streamersBucket)

	seq, _ := bucket.NextSequence()
	streamer.ID = int(seq)
	streamer.CreatedAt = time.Now()
	streamer.UpdatedAt = time.Now()

	buf, err := json.Marshal(streamer)
	if err != nil {
		return err
	}

	if err := bucket.Put(itob(streamer.ID), buf); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *boltStore) UpdateStreamer(streamer *models.Streamer) error {
	if streamer.ID <= 0 {
		return errors.New("invalid service id")
	}

	tx, err := s.Begin(true)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	bucket := tx.Bucket(streamersBucket)

	streamer.UpdatedAt = time.Now()

	buf, err := json.Marshal(streamer)
	if err != nil {
		return err
	}

	if err := bucket.Put(itob(streamer.ID), buf); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *boltStore) DeleteStreamer(id int) error {
	return s.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(streamersBucket).Delete(itob(id))
	})
}
