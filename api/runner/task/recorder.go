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

type createRecorderTask struct {
	recorder *api.RecorderConfig
}

func CreateRecorder(recorder *api.RecorderConfig) runner.Task {
	return &createRecorderTask{
		recorder: recorder,
	}
}

func (t *createRecorderTask) ID() runner.TaskID {
	return runner.TaskCreateRecorder
}

func (t *createRecorderTask) Run(ctx context.Context) (err error) {
	if t.recorder == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create recorder %s: %v", t.recorder.Name, err))
	}()

	v, err := json.Marshal(t.recorder)
	if err != nil {
		return err
	}

	if err := client.Default().CreateRecorder(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateRecorderTask struct {
	recorder *api.RecorderConfig
}

func UpdateRecorder(recorder *api.RecorderConfig) runner.Task {
	return &updateRecorderTask{
		recorder: recorder,
	}
}

func (t *updateRecorderTask) ID() runner.TaskID {
	return runner.TaskUpdateRecorder
}

func (t *updateRecorderTask) Run(ctx context.Context) (err error) {
	if t.recorder == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update recorder %s: %v", t.recorder.Name, err))
	}()

	v, err := json.Marshal(t.recorder)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateRecorder(ctx, t.recorder.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteRecorderTask struct {
	recorder string
}

func DeleteRecorder(recorder string) runner.Task {
	return &deleteRecorderTask{
		recorder: recorder,
	}
}

func (t *deleteRecorderTask) ID() runner.TaskID {
	return runner.TaskDeleteRecorder
}

func (t *deleteRecorderTask) Run(ctx context.Context) (err error) {
	if t.recorder == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete recorder %s: %v", t.recorder, err))
	}()

	if err := client.Default().DeleteRecorder(ctx, t.recorder); err != nil {
		return err
	}
	return nil
}
