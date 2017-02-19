package task

import (
	"github.com/sohlich/ticktock/config"
	"gopkg.in/mgo.v2/bson"
)

type TaskRepository interface {
	FindById(id string) (*Task, error)
	FindAllByOwner(ownerId string, limit int) ([]*Task, error)
	FindAllByStatusAndOwner(status string, ownerId string) ([]*Task, error)
	Save(t *Task) error
	InsertTags(id string, tags []string) error
}

type MgoTaskRepository struct {
	config.MgoRepository
}

var Repository TaskRepository

func InitializeRepository(db *config.Database) error {
	var err error
	Repository = &MgoTaskRepository{
		config.MgoRepository{
			db.Database().C("tasks"),
		},
	}
	return err
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

func (m *MgoTaskRepository) InsertTags(id string, tags []string) error {
	err := m.Update(bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": tags}}})
	return err
}
