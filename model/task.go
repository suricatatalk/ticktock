package model

import "github.com/goadesign/goa/uuid"
import "time"

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
