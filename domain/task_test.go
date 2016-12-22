package domain

import "testing"

func TestTransitions(t *testing.T) {
	instance := &Task{
		Status: "running",
	}

	instance.ChangeState(Pause)

	if instance.Status != StatusPaused {
		t.Error("Status change failed for Pause")
		t.Fail()
		return
	}

	instance.ChangeState(Start)

	if instance.Status != StatusRunning {
		t.Error("Status change failed for Start after Pause")
		t.Fail()
		return
	}

	instance.ChangeState(Finish)

	if instance.Status != StatusFinished {
		t.Error("Status change failed for Finish after Resume")
		t.Fail()
		return
	}

}
