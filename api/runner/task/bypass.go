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

type createBypassTask struct {
	bypass *api.BypassConfig
}

func CreateBypass(bypass *api.BypassConfig) runner.Task {
	return &createBypassTask{
		bypass: bypass,
	}
}

func (t *createBypassTask) ID() runner.TaskID {
	return runner.TaskCreateBypass
}

func (t *createBypassTask) Run(ctx context.Context) (err error) {
	if t.bypass == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create bypass %s: %v", t.bypass.Name, err))
	}()

	v, err := json.Marshal(t.bypass)
	if err != nil {
		return err
	}

	if err := client.Default().CreateBypass(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateBypassTask struct {
	bypass *api.BypassConfig
}

func UpdateBypass(bypass *api.BypassConfig) runner.Task {
	return &updateBypassTask{
		bypass: bypass,
	}
}

func (t *updateBypassTask) ID() runner.TaskID {
	return runner.TaskUpdateBypass
}

func (t *updateBypassTask) Run(ctx context.Context) (err error) {
	if t.bypass == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update bypass %s: %v", t.bypass.Name, err))
	}()

	v, err := json.Marshal(t.bypass)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateBypass(ctx, t.bypass.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteBypassTask struct {
	bypass string
}

func DeleteBypass(bypass string) runner.Task {
	return &deleteBypassTask{
		bypass: bypass,
	}
}

func (t *deleteBypassTask) ID() runner.TaskID {
	return runner.TaskDeleteBypass
}

func (t *deleteBypassTask) Run(ctx context.Context) (err error) {
	if t.bypass == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete bypass %s: %v", t.bypass, err))
	}()

	if err := client.Default().DeleteBypass(ctx, t.bypass); err != nil {
		return err
	}
	return nil
}
