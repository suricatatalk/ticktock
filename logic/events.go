package logic

import "github.com/sohlich/ticktock/security"
import "github.com/sohlich/ticktock/domain"
import "fmt"
import "time"

type EventDTO struct {
	TaskName        string `bson:"taskName" json:"taskName"`
	TaskID          string `bson:"taskId" json:"taskId"`
	EventEpoch      int64  `bson:"eventEpoch" json:"eventEpoch"`
	EventTypeString string `bson:"eventType" json:"eventType"`
}

type EventFunction func(user security.User, event *EventDTO) (*domain.Task, error)

func Start(user security.User, event *EventDTO) (*domain.Task, error) {
	tasks, err := domain.Tasks.FindAllByStatusAndOwner("running", user.ID)
	for _, val := range tasks {
		val.ChangeState(domain.Finish)
		domain.Tasks.Save(val)
	}

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
	return changeState(domain.Pause, user, event)
}

func Resume(user security.User, event *EventDTO) (*domain.Task, error) {
	return changeState(domain.Start, user, event)
}

func Finish(user security.User, event *EventDTO) (*domain.Task, error) {
	return changeState(domain.Finish, user, event)
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
	domain.Tasks.Save(task)
	return task, nil
}
