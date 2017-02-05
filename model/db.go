package model

import (
	"gopkg.in/mgo.v2"
)

var DB = Database{}

type Database struct {
	sess *mgo.Session
	db   *mgo.Database
}

func (db *Database) Session() *mgo.Session {
	return db.sess
}

func (db *Database) Database() *mgo.Database {
	return db.db
}

func (db *Database) Close() {
	db.sess.Close()
}

func (db *Database) Open(cfg StorageConfig) error {
	var err error
	db.sess, err = mgo.Dial(cfg.ConnectionString)
	if err != nil {
		return err
	}
	db.db = db.sess.DB(cfg.Database)
	return nil
}
