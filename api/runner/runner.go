package runner

import (
	"context"
	"sync"
	"time"
)

var (
	runner = NewRunner()
)

func Default() *Runner {
	return runner
}

type TaskEvent struct {
	TaskID string
	Err    error
}

type taskState struct {
	task   Task
	cancel context.CancelFunc
}

type Options struct {
	Async    bool
	Interval time.Duration
}

type Option func(opts *Options)

func WithAync(aync bool) Option {
	return func(opts *Options) {
		opts.Async = aync
	}
}

func WithInterval(interval time.Duration) Option {
	return func(opts *Options) {
		opts.Interval = interval
	}
}

type Runner struct {
	events chan *TaskEvent
	states map[string]taskState
	mu     sync.RWMutex
}

func NewRunner() *Runner {
	return &Runner{
		events: make(chan *TaskEvent, 16),
		states: make(map[string]taskState),
	}
}

func (r *Runner) Event() <-chan *TaskEvent {
	return r.events
}

func (r *Runner) Exec(ctx context.Context, task Task, opts ...Option) error {
	if task == nil || task.ID() == "" {
		return nil
	}

	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if !options.Async {
		return task.Run(ctx)
	}

	r.Cancel(task.ID())

	ctx, cancel := context.WithCancel(ctx)
	r.setState(taskState{
		task:   task,
		cancel: cancel,
	})

	go func() {
		defer cancel()

		run := func() {
			if err := task.Run(ctx); err != nil {
				r.events <- &TaskEvent{
					TaskID: task.ID(),
					Err:    err,
				}
				return
			}
			r.events <- &TaskEvent{
				TaskID: task.ID(),
			}
		}

		run()

		interval := options.Interval
		if interval <= 0 {
			return
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				run()
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *Runner) Cancel(id string) {
	r.delState(id)
}

func (r *Runner) setState(state taskState) {
	if state.task == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.states[state.task.ID()] = state
}

func (r *Runner) delState(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	state := r.states[id]
	if state.cancel != nil {
		state.cancel()
	}

	delete(r.states, id)
}
