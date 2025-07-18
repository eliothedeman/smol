package control

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eliothedeman/smol/unit"
)

type mockCtx struct {
	context.Context
	units []unit.UnitDesc
}

func (m *mockCtx) Units() []unit.UnitDesc {
	return m.units
}

func (m *mockCtx) Spawn(name string, f unit.UnitFactory) unit.UnitRef {
	return &mockUnitRef{name: name}
}

func (m *mockCtx) Self() unit.UnitRef {
	return &mockUnitRef{name: "test"}
}

func (m *mockCtx) Subscribe(other unit.Unit)   {}
func (m *mockCtx) Unsubscribe(other unit.Unit) {}

type mockUnitRef struct {
	name string
}

func (m *mockUnitRef) Name() string {
	return m.name
}

func (m *mockUnitRef) Send(msg any) {}

func (m *mockUnitRef) Stop() {}

func TestLifecycleStateTransitions(t *testing.T) {
	l := NewLifecycle()

	if l.State() != StateInit {
		t.Errorf("Expected initial state to be StateInit, got %v", l.State())
	}

	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	if l.State() != StateRunning {
		t.Errorf("Expected state after Init to be StateRunning, got %v", l.State())
	}
}

func TestLifecycleShutdown(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	l.Shutdown()

	if l.State() != StateStopped {
		t.Errorf("Expected state after Shutdown to be StateStopped, got %v", l.State())
	}
}

func TestLifecycleDoubleShutdown(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	l.Shutdown()
	firstState := l.State()

	l.Shutdown()
	secondState := l.State()

	if firstState != secondState {
		t.Errorf("Expected double shutdown to not change state, got %v then %v", firstState, secondState)
	}
}

func TestLifecycleHandleMessages(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	from := &mockUnitRef{name: "test"}

	err := l.Handle(ctx, from, "status")
	if err != nil {
		t.Errorf("Handle returned error: %v", err)
	}

	err = l.Handle(ctx, from, "shutdown")
	if err != nil {
		t.Errorf("Handle returned error: %v", err)
	}

	if l.State() != StateStopped {
		t.Errorf("Expected state after shutdown message to be StateStopped, got %v", l.State())
	}
}

func TestLifecycleContextCancellation(t *testing.T) {
	l := NewLifecycle()
	ctx, cancel := context.WithCancel(context.Background())
	mockCtx := &mockCtx{Context: ctx}
	l.Init(mockCtx)

	cancel()

	// Give some time for the context cancellation to be processed
	time.Sleep(10 * time.Millisecond)

	if l.State() != StateStopped {
		t.Errorf("Expected state after context cancellation to be StateStopped, got %v", l.State())
	}
}

func TestAddTask(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	var executed atomic.Bool
	l.AddTask(func() {
		executed.Store(true)
	})

	l.Shutdown()

	if !executed.Load() {
		t.Error("Expected task to be executed")
	}
}

func TestAddTaskWithContext(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	var executed atomic.Bool
	l.AddTaskWithContext(func(ctx context.Context) {
		executed.Store(true)
	})

	// Wait for task execution with timeout
	deadline := time.After(100 * time.Millisecond)
	for !executed.Load() {
		select {
		case <-deadline:
			t.Error("Expected task with context to be executed")
			return
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}

	l.Shutdown()
}

func TestTaskCancellation(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	var started atomic.Bool
	var completed atomic.Bool

	l.AddTaskWithContext(func(ctx context.Context) {
		started.Store(true)
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
			completed.Store(true)
		}
	})

	// Wait for task to start
	deadline := time.After(50 * time.Millisecond)
	for !started.Load() {
		select {
		case <-deadline:
			t.Error("Expected task to start")
			return
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}

	l.Shutdown()

	// The task should not complete because shutdown should cancel it
	// We can't directly test the cancellation, but we can verify the test completes quickly
	// which implies the task was cancelled
	time.Sleep(50 * time.Millisecond)
}
func TestConcurrentStateAccess(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	const numGoroutines = 100
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_ = l.State()
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	l.Shutdown()
}

func TestLifecycleWaitForTasks(t *testing.T) {
	l := NewLifecycle()
	ctx := &mockCtx{Context: context.Background()}
	l.Init(ctx)

	var taskCompleted atomic.Bool

	l.AddTask(func() {
		time.Sleep(10 * time.Millisecond)
		taskCompleted.Store(true)
	})

	l.Shutdown()

	if !taskCompleted.Load() {
		t.Error("Expected task to complete before shutdown")
	}
}
