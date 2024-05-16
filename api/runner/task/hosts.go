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

type createHostMapperTask struct {
	hostMapper *api.HostsConfig
}

func CreateHostMapper(hostMapper *api.HostsConfig) runner.Task {
	return &createHostMapperTask{
		hostMapper: hostMapper,
	}
}

func (t *createHostMapperTask) ID() runner.TaskID {
	return runner.TaskCreateHosts
}

func (t *createHostMapperTask) Run(ctx context.Context) (err error) {
	if t.hostMapper == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create hosts %s: %v", t.hostMapper.Name, err))
	}()

	v, err := json.Marshal(t.hostMapper)
	if err != nil {
		return err
	}

	if err := client.Default().CreateHostMapper(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateHostMapperTask struct {
	hostMapper *api.HostsConfig
}

func UpdateHostMapper(hostMapper *api.HostsConfig) runner.Task {
	return &updateHostMapperTask{
		hostMapper: hostMapper,
	}
}

func (t *updateHostMapperTask) ID() runner.TaskID {
	return runner.TaskUpdateHosts
}

func (t *updateHostMapperTask) Run(ctx context.Context) (err error) {
	if t.hostMapper == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update hosts %s: %v", t.hostMapper.Name, err))
	}()

	v, err := json.Marshal(t.hostMapper)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateHostMapper(ctx, t.hostMapper.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteHostMapperTask struct {
	hostMapper string
}

func DeleteHostMapper(hostMapper string) runner.Task {
	return &deleteHostMapperTask{
		hostMapper: hostMapper,
	}
}

func (t *deleteHostMapperTask) ID() runner.TaskID {
	return runner.TaskDeleteHosts
}

func (t *deleteHostMapperTask) Run(ctx context.Context) (err error) {
	if t.hostMapper == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete hosts %s: %v", t.hostMapper, err))
	}()

	if err := client.Default().DeleteHostMapper(ctx, t.hostMapper); err != nil {
		return err
	}
	return nil
}
