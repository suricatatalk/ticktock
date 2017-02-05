package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/sohlich/ticktock/model"
	"github.com/sohlich/ticktock/security"
)

type EventDTO struct {
	TaskName        string `bson:"taskName" json:"taskName"`
	TaskID          string `bson:"taskId" json:"taskId"`
	EventEpoch      int64  `bson:"eventEpoch" json:"eventEpoch"`
	EventTypeString string `bson:"eventType" json:"eventType"`
}

type EventFunction func(user security.User, event *EventDTO) (*model.Task, error)

func Start(user security.User, event *EventDTO) (*model.Task, error) {
	tasks, err := model.Tasks.FindAllByStatusAndOwner("running", user.ID)
	if len(tasks) > 0 {
		return nil, fmt.Errorf("Another task is running")
	}

	// for _, val := range tasks {
	// 	val.ChangeState(model.Finish)
	// 	model.Tasks.Save(val)
	// }

	task := model.NewTask()
	task.Status = model.StatusRunning
	task.OwnerID = user.ID
	task.Name = event.TaskName
	task.Start = time.Now().Unix()
	task.AddEvent(model.Start, task.Start)
	err = model.Tasks.Save(task)
	return task, err
}

func Pause(user security.User, event *EventDTO) (*model.Task, error) {
	logChange("Pause", user.ID, event.TaskID)
	return changeState(model.Pause, user, event)
}

func Resume(user security.User, event *EventDTO) (*model.Task, error) {
	logChange("Resume", user.ID, event.TaskID)
	return changeState(model.Start, user, event)
}

func Finish(user security.User, event *EventDTO) (*model.Task, error) {
	logChange("Finish", user.ID, event.TaskID)
	return changeState(model.Finish, user, event)
}
func logChange(event, userID, taskID string) {
	log.Printf("[%s] %s task_id: %s ", userID, event, taskID)
}

func changeState(action int, user security.User, event *EventDTO) (*model.Task, error) {
	task, err := model.Tasks.FindById(event.TaskID)
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
	model.Tasks.Save(task)
	return task, nil
}
