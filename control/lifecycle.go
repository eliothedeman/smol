package control

import (
	"context"
	"fmt"
	"sync"

	"github.com/eliothedeman/smol/unit"
)

type LifecycleState int

const (
	StateInit LifecycleState = iota
	StateStarting
	StateRunning
	StateStopping
	StateStopped
)

func (s LifecycleState) String() string {
	switch s {
	case StateInit:
		return "init"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type Lifecycle struct {
	mu       sync.RWMutex
	state    LifecycleState
	shutdown chan struct{}
	wg       sync.WaitGroup
}

func NewLifecycle() *Lifecycle {
	return &Lifecycle{
		state:    StateInit,
		shutdown: make(chan struct{}),
	}
}

func (l *Lifecycle) Init(ctx unit.Ctx) {
	l.mu.Lock()
	l.state = StateStarting
	l.mu.Unlock()

	go l.monitorShutdown(ctx)

	l.mu.Lock()
	l.state = StateRunning
	l.mu.Unlock()
}

func (l *Lifecycle) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	switch msg := message.(type) {
	case string:
		switch msg {
		case "status":
			l.mu.RLock()
			state := l.state
			l.mu.RUnlock()
			from.Send(fmt.Sprintf("lifecycle state: %s", state))
		case "shutdown":
			l.Shutdown()
		}
	case context.Context:
		select {
		case <-msg.Done():
			l.Shutdown()
		default:
		}
	}
	return nil
}

func (l *Lifecycle) State() LifecycleState {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.state
}

func (l *Lifecycle) Shutdown() {
	l.mu.Lock()
	if l.state == StateStopping || l.state == StateStopped {
		l.mu.Unlock()
		return
	}
	l.state = StateStopping
	close(l.shutdown)
	l.mu.Unlock()

	l.wg.Wait()

	l.mu.Lock()
	l.state = StateStopped
	l.mu.Unlock()
}

func (l *Lifecycle) monitorShutdown(ctx unit.Ctx) {
	select {
	case <-l.shutdown:
		return
	case <-ctx.Done():
		l.Shutdown()
		return
	}
}

func (l *Lifecycle) AddTask(fn func()) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		fn()
	}()
}

func (l *Lifecycle) AddTaskWithContext(fn func(ctx context.Context)) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		select {
		case <-l.shutdown:
			return
		default:
			fn(ctx)
		}
	}()
}
