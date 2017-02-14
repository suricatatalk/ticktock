package task

import "testing"

func TestTransitions(t *testing.T) {
	instance := &Task{
		Status: "running",
	}

	instance.ChangeState(EventPause)

	if instance.Status != StatusPaused {
		t.Error("Status change failed for Pause")
		t.Fail()
		return
	}

	instance.ChangeState(EventStart)

	if instance.Status != StatusRunning {
		t.Error("Status change failed for Start after Pause")
		t.Fail()
		return
	}

	instance.ChangeState(EventFinish)

	if instance.Status != StatusFinished {
		t.Error("Status change failed for Finish after Resume")
		t.Fail()
		return
	}

	instance.Status = StatusPaused

	instance.ChangeState(EventFinish)

	if instance.Status != StatusFinished {
		t.Error("Status change failed for Finish after Pause")
		t.Fail()
		return
	}

	instance.Status = StatusFinished

	instance.ChangeState(EventStart)

	if instance.Status != StatusFinished {
		t.Error("Status change failed for Pause")
		t.Fail()
		return
	}

}
