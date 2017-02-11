package model

import (
	"log"
	"testing"
)

func TestOpen(t *testing.T) {

	cfg := Environment{
		ConnectionString: "mongodb://localhost/ticktock_test",
	}

	err := Open(cfg)
	if err != nil {
		log.Println("Cannot open db session: " + err.Error())
		t.FailNow()
	}
	defer Close()

	task := NewTask()
	err = Tasks.Save(task)
	if err != nil {
		log.Println("Cannot save task: " + err.Error())
		t.FailNow()
	}
}
