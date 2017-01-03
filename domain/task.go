package domain

import "log"
import "github.com/goadesign/goa/uuid"
import "time"

// Events
const (
	Start  = iota
	Pause  = iota
	Finish = iota
)

const (
	StatusRunning  = "running"
	StatusPaused   = "paused"
	StatusFinished = "finished"
)

var (
	transitions map[string]map[int]string
)

func init() {
	transitions = make(map[string]map[int]string)

	// transitions for running state
	running := make(map[int]string)
	running[Pause] = StatusPaused
	running[Finish] = StatusFinished
	transitions[StatusRunning] = running

	// Transitions for paused state
	paused := make(map[int]string)
	paused[Finish] = StatusFinished
	paused[Start] = StatusRunning
	transitions[StatusPaused] = paused

	log.Println("Transitions initialized")
	log.Printf("%v", transitions)
}

type Event struct {
	EventEpoch int64 `bson:"eventEpoch" json:"eventEpoch"`
	EventType  int   `bson:"eventType" json:"eventType"`
	Duration   int64 `bson:"duration" json:"duration"`
}

type Task struct {
	ID       string  `bson:"_id" json:"id"`
	Name     string  `bson:"name" json:"name"`
	OwnerID  string  `bson:"ownerId" json:"ownerId"`
	Status   string  `bson:"status" json:"status"`
	Start    int64   `bson:"start" json:"start"`
	End      int64   `bson:"end" json:"end"`
	Duration int64   `bson:"duration" json:"duration"`
	Events   []Event `bson:"events" json:"events"`
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
	t.AddEvent(Start, start)
}

func (t *Task) ChangeState(transition int) {
	if next, ok := transitions[t.Status][transition]; ok {
		t.Status = next
		t.AddEvent(transition, time.Now().Unix())
	}
}

func (t *Task) AddEvent(event int, timestamp int64) {

	eCnt := len(t.Events)
	duration := int64(0)
	if (event == Pause || event == Finish) && eCnt > 0 {
		duration = timestamp - t.Events[eCnt-1].EventEpoch
	}

	t.Events = append(t.Events, Event{
		EventEpoch: timestamp,
		EventType:  event,
		Duration:   duration,
	})
}
