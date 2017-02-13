package model

import (
	"time"

	"github.com/goadesign/goa/uuid"
	"gopkg.in/mgo.v2/bson"
)

// Events

const (
	StatusRunning  = "running"
	StatusPaused   = "paused"
	StatusFinished = "finished"
)

type Task struct {
	ID       string   `bson:"_id" json:"id"`
	Name     string   `bson:"name" json:"name"`
	OwnerID  string   `bson:"ownerId" json:"ownerId"`
	Status   string   `bson:"status" json:"status"`
	Start    int64    `bson:"start" json:"start"`
	End      int64    `bson:"end" json:"end"`
	Duration int64    `bson:"duration" json:"duration"`
	Events   []Event  `bson:"events" json:"events"`
	Tags     []string `bson:"tags" json:"tags"`
}

func NewTask() *Task {
	return &Task{
		ID: uuid.NewV4().String(),
	}
}

func (t *Task) GenerateID() {
	t.ID = uuid.NewV4().String()
}

func (t *Task) StartTask() {
	start := time.Now().Unix()
	t.Start = start
	t.Status = StatusRunning
	t.AddEvent(EventStart, start)
}

func (t *Task) ChangeState(transition string) {
	if next, ok := transitions[t.Status][transition]; ok {
		t.Status = next
		if transition == EventFinish {
			t.End = time.Now().Unix()
		}
		t.AddEvent(transition, time.Now().Unix())
	}
}

func (t *Task) AddEvent(event string, timestamp int64) {

	eCnt := len(t.Events)
	duration := int64(0)
	if (event == EventPause || event == EventFinish) && eCnt > 0 {
		duration = timestamp - t.Events[eCnt-1].EventEpoch
	}

	t.Events = append(t.Events, Event{
		EventEpoch: timestamp,
		EventType:  event,
		Duration:   duration,
	})
}

type TaskRepository interface {
	FindById(id string) (*Task, error)
	FindAllByOwner(ownerId string, limit int) ([]*Task, error)
	FindAllByStatusAndOwner(status string, ownerId string) ([]*Task, error)
	Save(t *Task) error
	InsertTags(id string, tags []string) error
}

type MgoTaskRepository struct {
	MgoRepository
}

var Tasks TaskRepository

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
