package logic

import "github.com/sohlich/ticktock/security"
import "github.com/sohlich/ticktock/domain"

type EventDTO struct {
	TaskName        string `bson:"taskName" json:"taskName"`
	TaskID          string `bson:"taskId" json:"taskId"`
	EventEpoch      int64  `bson:"eventEpoch" json:"eventEpoch"`
	EventTypeString string `bson:"eventType" json:"eventType"`
}

func Start(user security.User, event *EventDTO) (error, *domain.Task) {
	tasks, err := domain.Tasks.FindAllByStatusAndOwner("running", user.ID)
	for _, val := range tasks {
		val.ChangeState(domain.Finish)
		domain.Tasks.Save(val)
	}

	task := domain.NewTask()
	task.Status = domain.StatusRunning
	task.OwnerID = user.ID
	task.Name = event.TaskName
	task.Start = event.EventEpoch
	task.AddEvent(domain.Start, event.EventEpoch)
	err = domain.Tasks.Save(task)
	return err, task
}
