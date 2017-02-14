package task

import (
	"log"
	"testing"

	"github.com/sohlich/ticktock/config"
)

func TestOpen(t *testing.T) {

	cfg := Environment{
		ConnectionString: "mongodb://localhost/ticktock_test",
	}

	err := config.DB.Open(cfg)
	InitializeRepository(DB)
	if err != nil {
		log.Println("Cannot open db session: " + err.Error())
		t.FailNow()
	}
	defer DB.Close()

	task := NewTask()
	err = Tasks.Save(task)
	if err != nil {
		log.Println("Cannot save task: " + err.Error())
		t.FailNow()
	}
}
