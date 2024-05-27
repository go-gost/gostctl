package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/api/client"
	"github.com/go-gost/gostctl/api/runner"
)

type createObserverTask struct {
	observer *api.ObserverConfig
}

func CreateObserver(observer *api.ObserverConfig) runner.Task {
	return &createObserverTask{
		observer: observer,
	}
}

func (t *createObserverTask) ID() runner.TaskID {
	return runner.TaskCreateObserver
}

func (t *createObserverTask) Run(ctx context.Context) (err error) {
	if t.observer == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create observer %s: %v", t.observer.Name, err))
	}()

	v, err := json.Marshal(t.observer)
	if err != nil {
		return err
	}

	if err := client.Default().CreateObserver(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateObserverTask struct {
	observer *api.ObserverConfig
}

func UpdateObserver(observer *api.ObserverConfig) runner.Task {
	return &updateObserverTask{
		observer: observer,
	}
}

func (t *updateObserverTask) ID() runner.TaskID {
	return runner.TaskUpdateObserver
}

func (t *updateObserverTask) Run(ctx context.Context) (err error) {
	if t.observer == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update observer %s: %v", t.observer.Name, err))
	}()

	v, err := json.Marshal(t.observer)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateObserver(ctx, t.observer.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteObserverTask struct {
	observer string
}

func DeleteObserver(observer string) runner.Task {
	return &deleteObserverTask{
		observer: observer,
	}
}

func (t *deleteObserverTask) ID() runner.TaskID {
	return runner.TaskDeleteObserver
}

func (t *deleteObserverTask) Run(ctx context.Context) (err error) {
	if t.observer == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete observer %s: %v", t.observer, err))
	}()

	if err := client.Default().DeleteObserver(ctx, t.observer); err != nil {
		return err
	}
	return nil
}
