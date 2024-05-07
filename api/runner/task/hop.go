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

type createHopTask struct {
	hop *api.HopConfig
}

func CreateHop(hop *api.HopConfig) runner.Task {
	return &createHopTask{
		hop: hop,
	}
}

func (t *createHopTask) ID() runner.TaskID {
	return runner.TaskCreateHop
}

func (t *createHopTask) Run(ctx context.Context) (err error) {
	if t.hop == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create hop %s: %v", t.hop.Name, err))
	}()

	v, err := json.Marshal(t.hop)
	if err != nil {
		return err
	}

	if err := client.Default().CreateHop(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateHopTask struct {
	hop *api.HopConfig
}

func UpdateHop(hop *api.HopConfig) runner.Task {
	return &updateHopTask{
		hop: hop,
	}
}

func (t *updateHopTask) ID() runner.TaskID {
	return runner.TaskUpdateHop
}

func (t *updateHopTask) Run(ctx context.Context) (err error) {
	if t.hop == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update hop %s: %v", t.hop.Name, err))
	}()

	v, err := json.Marshal(t.hop)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateHop(ctx, t.hop.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteHopTask struct {
	hop string
}

func DeleteHop(hop string) runner.Task {
	return &deleteHopTask{
		hop: hop,
	}
}

func (t *deleteHopTask) ID() runner.TaskID {
	return runner.TaskDeleteHop
}

func (t *deleteHopTask) Run(ctx context.Context) (err error) {
	if t.hop == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete hop %s: %v", t.hop, err))
	}()

	if err := client.Default().DeleteHop(ctx, t.hop); err != nil {
		return err
	}
	return nil
}
