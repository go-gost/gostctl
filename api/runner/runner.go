package runner

import (
	"context"
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

type Runner struct {
	events chan *TaskEvent
}

func NewRunner() *Runner {
	return &Runner{
		events: make(chan *TaskEvent, 16),
	}
}

func (r *Runner) Event() <-chan *TaskEvent {
	return r.events
}

func (r *Runner) Exec(ctx context.Context, task Task) error {
	if task == nil {
		return nil
	}

	return task.Run(ctx)
}

func (r *Runner) ExecAsync(ctx context.Context, task Task, interval time.Duration) error {
	if task == nil {
		return nil
	}

	go func() {
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
