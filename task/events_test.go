package task_test

import (
	"testing"

	"fmt"

	"github.com/sohlich/ticktock/task"
)

type NoOpTaskRepo struct {
	CalledID     string
	OwnerID      string
	Status       string
	Task         *task.Task
	FindFunction FindFunc
}

type FindFunc func(id string) (*task.Task, error)

func (r *NoOpTaskRepo) Reset() {
	r.CalledID = ""
	r.OwnerID = ""
	r.Status = ""
	r.Task = nil
}

func (r *NoOpTaskRepo) FindById(id string) (*task.Task, error) {

	if r.FindFunction != nil {
		return r.FindFunction(id)
	}

	return r.Task, nil
}

func (r *NoOpTaskRepo) FindAllByOwner(ownerId string, limit int) ([]*task.Task, error) {
	return nil, nil
}

func (r *NoOpTaskRepo) FindAllByStatusAndOwner(status string, ownerId string) ([]*task.Task, error) {
	r.OwnerID = ownerId
	r.Status = status
	return nil, nil
}

func (r *NoOpTaskRepo) Save(t *task.Task) error {
	return nil
}

func (r *NoOpTaskRepo) InsertTags(id string, tags []string) error {
	return nil
}

func TestStart(t *testing.T) {
	mock := &NoOpTaskRepo{}
	task.Tasks = mock
	event := task.Event{
		TaskName:  "Test task",
		EventType: "start",
	}
	testUser := task.User{
		ID: "1234@test",
	}
	task, _ := task.Start(testUser, &event)

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
	task.Tasks = mock
	event := &task.Event{
		TaskName:  "Test task",
		EventType: "start",
	}
	testUser := task.User{
		ID: "1234@test",
	}
	task, _ := task.Start(testUser, event)

	event.TaskID = "1234"
	event.EventType = "pause"

	mock.FindFunction = func(id string) (*task.Task, error) {
		if id == event.TaskID {
			return task, nil
		}
		return nil, fmt.Errorf("Cannot find record with id:%s\n", id)
	}

	var err error
	task, err = task.Pause(testUser, event)

	if err != nil {
		t.Error("ID of task does not match")
		t.FailNow()
	}

	if len(task.Events) != 2 {
		t.Error("Not enough events")
		t.FailNow()
	}

	if task.Events[1].EventType != task.EventPause {
		t.Errorf("Bad second event %v:", task.Events[1])
		t.FailNow()
	}

}
