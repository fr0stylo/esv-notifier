package database

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"

	"github.com/boltdb/bolt"
)

var (
	db    *bolt.DB
	once  sync.Once
	mutex sync.Mutex
)

type Database[T Entity] struct {
	db     *bolt.DB
	bucket string
}

func NewDatabase[T Entity](path string, bucket string) (*Database[T], error) {
	once.Do(func() {
		bdb, err := bolt.Open(path, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		db = bdb
	})

	return &Database[T]{
		db:     db,
		bucket: bucket,
	}, nil
}

func (d *Database[T]) Close() {
	d.db.Close()
}

type Entity interface {
	Insert(any) ([]byte, error)
	Update(any, []byte) ([]byte, error)
	Marshal() ([]byte, error)
}

func (d *Database[T]) Upsert(key string, value any) error {
	mutex.Lock()
	defer mutex.Unlock()

	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(d.bucket))
		if err != nil {
			return err
		}

		items := b.Get([]byte(key))
		var existing T
		if err := gob.NewDecoder(bytes.NewReader(items)).Decode(&existing); err != nil {
			return err
		}
		if items == nil {
			x, err := existing.Insert(value)
			if err != nil {
				return err
			}

			return b.Put([]byte(key), x)
		} else {
			x, err := existing.Update(value, items)
			if err != nil {
				return err
			}

			return b.Put([]byte(key), x)
		}
	})
}

func (d *Database[T]) Iterate(fn func(key string, value T) error) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucket))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var r T
			gob.NewDecoder(bytes.NewReader(v)).Decode(&r)
			return fn(string(k), r)
		})
	})
}

func (d *Database[T]) BatchUpdate(key [][]byte, fn func(value T) ([]byte, error)) error {
	mutex.Lock()
	defer mutex.Unlock()
	return d.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(d.bucket))
		if err != nil {
			return err
		}

		for _, k := range key {
			v := b.Get(k)

			var r T
			gob.NewDecoder(bytes.NewReader(v)).Decode(&r)
			i, err := fn(r)
			if err != nil {
				return err
			}

			if i == nil {
				b.Delete(k)
			} else {
				b.Put(k, i)
			}
		}

		return nil
	})
}

func (d *Database[T]) BatchInsert(items []T, keyFn func(T) ([]byte, error)) error {
	mutex.Lock()
	defer mutex.Unlock()
	return d.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(d.bucket))
		if err != nil {
			return err
		}

		for _, item := range items {
			key, err := keyFn(item)
			if err != nil {
				return err
			}

			buf, err := item.Marshal()
			if err != nil {
				return err
			}

			if err := b.Put(key, buf); err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *Database[T]) GetByPrefix(prefix string) ([]T, error) {
	items := []T{}
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {
			var r T
			gob.NewDecoder(bytes.NewReader(v)).Decode(&r)
			items = append(items, r)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	return items, nil
}
