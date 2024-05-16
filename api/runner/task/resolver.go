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

type createResolverTask struct {
	resolver *api.ResolverConfig
}

func CreateResolver(resolver *api.ResolverConfig) runner.Task {
	return &createResolverTask{
		resolver: resolver,
	}
}

func (t *createResolverTask) ID() runner.TaskID {
	return runner.TaskCreateResolver
}

func (t *createResolverTask) Run(ctx context.Context) (err error) {
	if t.resolver == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create resolver %s: %v", t.resolver.Name, err))
	}()

	v, err := json.Marshal(t.resolver)
	if err != nil {
		return err
	}

	if err := client.Default().CreateResolver(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateResolverTask struct {
	resolver *api.ResolverConfig
}

func UpdateResolver(resolver *api.ResolverConfig) runner.Task {
	return &updateResolverTask{
		resolver: resolver,
	}
}

func (t *updateResolverTask) ID() runner.TaskID {
	return runner.TaskUpdateResolver
}

func (t *updateResolverTask) Run(ctx context.Context) (err error) {
	if t.resolver == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update resolver %s: %v", t.resolver.Name, err))
	}()

	v, err := json.Marshal(t.resolver)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateResolver(ctx, t.resolver.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteResolverTask struct {
	resolver string
}

func DeleteResolver(resolver string) runner.Task {
	return &deleteResolverTask{
		resolver: resolver,
	}
}

func (t *deleteResolverTask) ID() runner.TaskID {
	return runner.TaskDeleteResolver
}

func (t *deleteResolverTask) Run(ctx context.Context) (err error) {
	if t.resolver == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete resolver %s: %v", t.resolver, err))
	}()

	if err := client.Default().DeleteResolver(ctx, t.resolver); err != nil {
		return err
	}
	return nil
}
