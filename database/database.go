package database

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// DB type
type DB struct {
	d *bolt.DB
}

// Init initialise DB type
func (d *DB) Init(dbname string) error {
	base, err := bolt.Open(dbname, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.Printf("database.Init: %s\n", err.Error())
		return err
	}
	d.d = base
	return nil
}

// GetConfig take config as *map[string]string from database
func (d *DB) GetConfig() *map[string]string {
	var conf map[string]string
	d.d.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("config"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			conf[string(k)] = string(v)
		}
		return nil
	})
	return &conf
}

// SetConfig save config in database
func (d *DB) SetConfig(conf *map[string]string) {
	d.d.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("config"))
		for k, v := range *conf {
			err := b.Put([]byte(k), []byte(v))
			if err != nil {
				log.Printf("database.SetConfig: %s\n", err.Error())
			}
		}
		return nil
	})
}
