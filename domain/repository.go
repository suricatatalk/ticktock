package domain

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
	FindAllByOwner(ownerId string) ([]*Task, error)
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
var session *mgo.Session

func Open(cfg StorageConfig) error {
	var err error
	session, err = mgo.Dial(cfg.ConnectionString)
	if err != nil {
		return err
	}
	Tasks = &MgoTaskRepository{
		MgoRepository{
			session.DB(cfg.Database).C("tasks"),
		},
	}
	return err
}

func Close() {
	if session != nil {
		session.Close()
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

func (m *MgoTaskRepository) FindAllByOwner(ownerId string) ([]*Task, error) {
	all := make([]*Task, 0)
	err := m.Find(bson.M{"ownerId": ownerId}).Sort("-start").All(&all)
	return all, err
}

func (m *MgoTaskRepository) FindAllByStatusAndOwner(status string, ownerId string) ([]*Task, error) {
	all := make([]*Task, 0)
	err := m.Find(bson.M{"ownerId": ownerId, "status": status}).Sort("-start").All(&all)
	return all, err
}
