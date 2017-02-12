package logic_test

import (
	"testing"

	"fmt"

	"github.com/sohlich/ticktock/logic"
	"github.com/sohlich/ticktock/model"
)

type NoOpTaskRepo struct {
	CalledID     string
	OwnerID      string
	Status       string
	Task         *model.Task
	FindFunction FindFunc
}

type FindFunc func(id string) (*model.Task, error)

func (r *NoOpTaskRepo) Reset() {
	r.CalledID = ""
	r.OwnerID = ""
	r.Status = ""
	r.Task = nil
}

func (r *NoOpTaskRepo) FindById(id string) (*model.Task, error) {

	if r.FindFunction != nil {
		return r.FindFunction(id)
	}

	return r.Task, nil
}

func (r *NoOpTaskRepo) FindAllByOwner(ownerId string, limit int) ([]*model.Task, error) {
	return nil, nil
}

func (r *NoOpTaskRepo) FindAllByStatusAndOwner(status string, ownerId string) ([]*model.Task, error) {
	r.OwnerID = ownerId
	r.Status = status
	return nil, nil
}

func (r *NoOpTaskRepo) Save(t *model.Task) error {
	return nil
}

func TestStart(t *testing.T) {
	mock := &NoOpTaskRepo{}
	model.Tasks = mock
	event := logic.EventDTO{
		TaskName:        "Test task",
		EventTypeString: "start",
	}
	testUser := model.User{
		ID: "1234@test",
	}
	task, _ := logic.Start(testUser, &event)

	if mock.OwnerID != testUser.ID {
		t.Errorf("OwnerID does not match")
		t.FailNow()
	}

	if task.Name != event.TaskName {
		t.Errorf("Name does not match")
		t.FailNow()
	}
}

func TestPause(t *testing.T) {
	mock := &NoOpTaskRepo{}
	model.Tasks = mock
	event := &logic.EventDTO{
		TaskName:        "Test task",
		EventTypeString: "start",
	}
	testUser := model.User{
		ID: "1234@test",
	}
	task, _ := logic.Start(testUser, event)

	event.TaskID = "1234"
	event.EventTypeString = "pause"

	mock.FindFunction = func(id string) (*model.Task, error) {
		if id == event.TaskID {
			return task, nil
		}
		return nil, fmt.Errorf("Cannot find record with id:%s\n", id)
	}

	var err error
	task, err = logic.Pause(testUser, event)

	if err != nil {
		t.Error("ID of task does not match")
		t.FailNow()
	}

	if len(task.Events) != 2 {
		t.Error("Not enough events")
		t.FailNow()
	}

	if task.Events[1].EventType != model.Pause {
		t.Errorf("Bad second event %v:", task.Events[1])
		t.FailNow()
	}

}
