package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/sohlich/ticktock/domain"
	"github.com/sohlich/ticktock/security"
)

type EventDTO struct {
	TaskName        string `bson:"taskName" json:"taskName"`
	TaskID          string `bson:"taskId" json:"taskId"`
	EventEpoch      int64  `bson:"eventEpoch" json:"eventEpoch"`
	EventTypeString string `bson:"eventType" json:"eventType"`
}

type EventFunction func(user security.User, event *EventDTO) (*domain.Task, error)

func Start(user security.User, event *EventDTO) (*domain.Task, error) {
	tasks, err := domain.Tasks.FindAllByStatusAndOwner("running", user.ID)
	if len(tasks) > 0 {
		return nil, fmt.Errorf("Another task is running")
	}

	// for _, val := range tasks {
	// 	val.ChangeState(domain.Finish)
	// 	domain.Tasks.Save(val)
	// }

	task := domain.NewTask()
	task.Status = domain.StatusRunning
	task.OwnerID = user.ID
	task.Name = event.TaskName
	task.Start = time.Now().Unix()
	task.AddEvent(domain.Start, task.Start)
	err = domain.Tasks.Save(task)
	return task, err
}

func Pause(user security.User, event *EventDTO) (*domain.Task, error) {
	logChange("Pause", user.ID, event.TaskID)
	return changeState(domain.Pause, user, event)
}

func Resume(user security.User, event *EventDTO) (*domain.Task, error) {
	logChange("Resume", user.ID, event.TaskID)
	return changeState(domain.Start, user, event)
}

func Finish(user security.User, event *EventDTO) (*domain.Task, error) {
	logChange("Finish", user.ID, event.TaskID)
	return changeState(domain.Finish, user, event)
}
func logChange(event, userID, taskID string) {
	log.Printf("[%s] %s task_id: %s ", userID, event, taskID)
}

func changeState(action int, user security.User, event *EventDTO) (*domain.Task, error) {
	task, err := domain.Tasks.FindById(event.TaskID)
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
	domain.Tasks.Save(task)
	return task, nil
}
