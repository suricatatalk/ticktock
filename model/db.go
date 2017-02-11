package model

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"gopkg.in/mgo.v2"
)

func init() {
	domainCfg := StorageConfig{}
	b, err := ioutil.ReadFile("database.json")
	if err != nil {
		log.Println(err.Error())
	}
	enc := json.NewDecoder(bytes.NewReader(b))
	err = enc.Decode(&domainCfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	DB.Open(domainCfg["Development"])
	InitializeRepository(DB)
}

type StorageConfig map[string]Environment

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
