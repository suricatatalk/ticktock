package logic

import (
	"testing"

	"fmt"

	"github.com/sohlich/ticktock/domain"
	"github.com/sohlich/ticktock/security"
)

type NoOpTaskRepo struct {
	CalledID     string
	OwnerID      string
	Status       string
	Task         *domain.Task
	FindFunction FindFunc
}

type FindFunc func(id string) (*domain.Task, error)

func (r *NoOpTaskRepo) Reset() {
	r.CalledID = ""
	r.OwnerID = ""
	r.Status = ""
	r.Task = nil
}

func (r *NoOpTaskRepo) FindById(id string) (*domain.Task, error) {

	if r.FindFunction != nil {
		return r.FindFunction(id)
	}

	return r.Task, nil
}

func (r *NoOpTaskRepo) FindAllByOwner(ownerId string) ([]*domain.Task, error) {
	return nil, nil
}

func (r *NoOpTaskRepo) FindAllByStatusAndOwner(status string, ownerId string) ([]*domain.Task, error) {
	r.OwnerID = ownerId
	r.Status = status
	return nil, nil
}

func (r *NoOpTaskRepo) Save(t *domain.Task) error {
	return nil
}

func TestStart(t *testing.T) {
	mock := &NoOpTaskRepo{}
	domain.Tasks = mock
	event := EventDTO{
		TaskName:        "Test task",
		EventTypeString: "start",
	}
	testUser := security.User{
		ID: "1234@test",
	}
	task, _ := Start(testUser, &event)

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
	domain.Tasks = mock
	event := &EventDTO{
		TaskName:        "Test task",
		EventTypeString: "start",
	}
	testUser := security.User{
		ID: "1234@test",
	}
	task, _ := Start(testUser, event)

	event.TaskID = "1234"
	event.EventTypeString = "pause"

	mock.FindFunction = func(id string) (*domain.Task, error) {
		if id == event.TaskID {
			return task, nil
		}
		return nil, fmt.Errorf("Cannot find record with id:%s\n", id)
	}

	var err error
	task, err = Pause(testUser, event)

	if err != nil {
		t.Error("ID of task does not match")
		t.FailNow()
	}

	if len(task.Events) != 2 {
		t.Error("Not enough events")
		t.FailNow()
	}

	if task.Events[1].EventType != domain.Pause {
		t.Errorf("Bad second event %v:", task.Events[1])
		t.FailNow()
	}

}
