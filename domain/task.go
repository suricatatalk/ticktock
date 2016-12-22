package domain

import "log"
import "github.com/goadesign/goa/uuid"

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
	EventEpoch int64 `json:"eventExpoch"`
	EventType  int   `json:"eventType"`
}

type Task struct {
	ID      string  `json:"_id"`
	OwnerID string  `json:"ownerId"`
	Status  string  `json:"status"`
	Start   int64   `json:"start"`
	End     int64   `json:"end"`
	Events  []Event `json:"events"`
}

func NewTask() *Task {
	return &Task{
		ID: uuid.NewV4().String(),
	}
}

func (t *Task) ChangeState(transition int) {
	if nextState, ok := transitions[t.Status][transition]; ok {
		log.Printf("Changing state to %s", nextState)
		t.Status = nextState
	}
}

func (t *Task) addEvent() {
	t.Events
}
