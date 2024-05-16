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

type createLimiterTask struct {
	limiter *api.LimiterConfig
}

func CreateLimiter(limiter *api.LimiterConfig) runner.Task {
	return &createLimiterTask{
		limiter: limiter,
	}
}

func (t *createLimiterTask) ID() runner.TaskID {
	return runner.TaskCreateLimiter
}

func (t *createLimiterTask) Run(ctx context.Context) (err error) {
	if t.limiter == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create limiter %s: %v", t.limiter.Name, err))
	}()

	v, err := json.Marshal(t.limiter)
	if err != nil {
		return err
	}

	if err := client.Default().CreateLimiter(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateLimiterTask struct {
	limiter *api.LimiterConfig
}

func UpdateLimiter(limiter *api.LimiterConfig) runner.Task {
	return &updateLimiterTask{
		limiter: limiter,
	}
}

func (t *updateLimiterTask) ID() runner.TaskID {
	return runner.TaskUpdateLimiter
}

func (t *updateLimiterTask) Run(ctx context.Context) (err error) {
	if t.limiter == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update limiter %s: %v", t.limiter.Name, err))
	}()

	v, err := json.Marshal(t.limiter)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateLimiter(ctx, t.limiter.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteLimiterTask struct {
	limiter string
}

func DeleteLimiter(limiter string) runner.Task {
	return &deleteLimiterTask{
		limiter: limiter,
	}
}

func (t *deleteLimiterTask) ID() runner.TaskID {
	return runner.TaskDeleteLimiter
}

func (t *deleteLimiterTask) Run(ctx context.Context) (err error) {
	if t.limiter == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete limiter %s: %v", t.limiter, err))
	}()

	if err := client.Default().DeleteLimiter(ctx, t.limiter); err != nil {
		return err
	}
	return nil
}
