package task

import (
	"fmt"
	"log"
	"time"

	"github.com/sohlich/ticktock/user"
)

const (
	EventStart  = "start"
	EventPause  = "pause"
	EventFinish = "finish"
)

var (
	transitions map[string]map[string]string
)

func init() {
	transitions = make(map[string]map[string]string)

	// transitions for running state
	running := make(map[string]string)
	running[EventPause] = StatusPaused
	running[EventFinish] = StatusFinished
	transitions[StatusRunning] = running

	// Transitions for paused state
	paused := make(map[string]string)
	paused[EventFinish] = StatusFinished
	paused[EventStart] = StatusRunning
	transitions[StatusPaused] = paused

	log.Println("Transitions initialized")
	log.Printf("%v", transitions)
}

type Event struct {
	TaskName   string `bson:"-" json:"taskName"`
	TaskID     string `bson:"-" json:"taskId"`
	EventEpoch int64  `bson:"eventEpoch" json:"eventEpoch"`
	EventType  string `bson:"eventType" json:"eventType"`
	Duration   int64  `bson:"duration" json:"duration"`
}

type EventFunction func(user user.User, event *Event) (*Task, error)

func Start(user user.User, event *Event) (*Task, error) {
	tasks, err := Tasks.FindAllByStatusAndOwner("running", user.ID)
	if len(tasks) > 0 {
		return nil, fmt.Errorf("Another task is running")
	}

	// for _, val := range tasks {
	// 	val.ChangeState(Finish)
	// 	Tasks.Save(val)
	// }

	task := NewTask()
	task.Status = StatusRunning
	task.OwnerID = user.ID
	task.Name = event.TaskName
	task.Start = time.Now().Unix()
	task.AddEvent(EventStart, task.Start)
	err = Tasks.Save(task)
	return task, err
}

func Pause(user user.User, event *Event) (*Task, error) {
	logChange("Pause", user.ID, event.TaskID)
	return changeState(EventPause, user, event)
}

func Resume(user user.User, event *Event) (*Task, error) {
	logChange("Resume", user.ID, event.TaskID)
	return changeState(EventStart, user, event)
}

func Finish(user user.User, event *Event) (*Task, error) {
	logChange("Finish", user.ID, event.TaskID)
	return changeState(EventFinish, user, event)
}
func logChange(event, userID, taskID string) {
	log.Printf("[%s] %s task_id: %s ", userID, event, taskID)
}

func changeState(action string, user user.User, event *Event) (*Task, error) {
	task, err := Tasks.FindById(event.TaskID)
	if err != nil {
		return nil, err
	}
	if task.OwnerID != user.ID {
		return nil, fmt.Errorf("User is not owner of task")
	}
	task.ChangeState(action)

	duration := int64(0)
	for _, tsk := range task.Events {
		duration = duration + tsk.Duration
	}
	task.Duration = duration
	Tasks.Save(task)
	return task, nil
}
