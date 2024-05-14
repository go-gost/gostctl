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

type createAutherTask struct {
	auther *api.AutherConfig
}

func CreateAuther(auther *api.AutherConfig) runner.Task {
	return &createAutherTask{
		auther: auther,
	}
}

func (t *createAutherTask) ID() runner.TaskID {
	return runner.TaskCreateAuther
}

func (t *createAutherTask) Run(ctx context.Context) (err error) {
	if t.auther == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create auther %s: %v", t.auther.Name, err))
	}()

	v, err := json.Marshal(t.auther)
	if err != nil {
		return err
	}

	if err := client.Default().CreateAuther(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateAutherTask struct {
	auther *api.AutherConfig
}

func UpdateAuther(auther *api.AutherConfig) runner.Task {
	return &updateAutherTask{
		auther: auther,
	}
}

func (t *updateAutherTask) ID() runner.TaskID {
	return runner.TaskUpdateAuther
}

func (t *updateAutherTask) Run(ctx context.Context) (err error) {
	if t.auther == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update auther %s: %v", t.auther.Name, err))
	}()

	v, err := json.Marshal(t.auther)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateAuther(ctx, t.auther.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteAutherTask struct {
	auther string
}

func DeleteAuther(auther string) runner.Task {
	return &deleteAutherTask{
		auther: auther,
	}
}

func (t *deleteAutherTask) ID() runner.TaskID {
	return runner.TaskDeleteAuther
}

func (t *deleteAutherTask) Run(ctx context.Context) (err error) {
	if t.auther == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete auther %s: %v", t.auther, err))
	}()

	if err := client.Default().DeleteAuther(ctx, t.auther); err != nil {
		return err
	}
	return nil
}
