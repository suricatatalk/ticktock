package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type StorageConfig struct {
	ConnectionString string
	Password         string
	Database         string
}

type TaskRepository interface {
	FindById(id string) (*Task, error)
	FindAllByOwner(ownerId string, limit int) ([]*Task, error)
	FindAllByStatusAndOwner(status string, ownerId string) ([]*Task, error)
	Save(t *Task) error
}

type MgoRepository struct {
	*mgo.Collection
}

type MgoTaskRepository struct {
	MgoRepository
}

var Tasks TaskRepository

// var session *mgo.Session

func Open(cfg StorageConfig) error {
	var err error
	DB.Open(cfg)
	Tasks = &MgoTaskRepository{
		MgoRepository{
			DB.Database().C("tasks"),
		},
	}
	return err
}

func Close() {
	if DB.Session() != nil {
		DB.Close()
	}
}

func (m *MgoTaskRepository) Save(t *Task) error {
	_, err := m.Upsert(bson.M{"_id": t.ID}, t)
	return err
}

func (m *MgoTaskRepository) FindById(id string) (*Task, error) {
	task := &Task{}
	err := m.Find(bson.M{"_id": id}).One(task)
	return task, err
}

func (m *MgoTaskRepository) FindAllByOwner(ownerId string, limit int) ([]*Task, error) {
	all := make([]*Task, 0)
	query := m.Find(bson.M{"ownerId": ownerId}).Sort("-start")
	if limit < 0 {
		query.Limit(limit)
	}

	err := query.All(&all)
	return all, err
}

func (m *MgoTaskRepository) FindAllByStatusAndOwner(status string, ownerId string) ([]*Task, error) {
	all := make([]*Task, 0)
	err := m.Find(bson.M{"ownerId": ownerId, "status": status}).Sort("-start").All(&all)
	return all, err
}

type UserRepository interface {
	FindById(id string) (*Task, error)
	Save(t *Task) error
}
