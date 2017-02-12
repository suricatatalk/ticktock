package model

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"gopkg.in/mgo.v2"
)

// Environment describes the configuration
// for environment, where the application
// runs
type Environment struct {
	ConnectionString string
	Password         string
	Database         string
}

// DB is the global object
// to hold session.
var DB = &Database{}

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

func (db *Database) Open(cfg Environment) error {
	var err error
	db.sess, err = mgo.Dial(cfg.ConnectionString)
	if err != nil {
		return err
	}
	db.db = db.sess.DB(cfg.Database)
	return nil
}

type MgoRepository struct {
	*mgo.Collection
}

func InitDB(cfgFile, env string) {
	var c map[string]Environment
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	enc := json.NewDecoder(bytes.NewReader(b))
	enc.Decode(c)
	if err != nil {
		log.Fatalln(err.Error())
	}
	DB.Open(c[env])
	InitializeRepository(DB)
}
