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

type createServiceTask struct {
	service *api.ServiceConfig
}

func CreateService(service *api.ServiceConfig) runner.Task {
	return &createServiceTask{
		service: service,
	}
}

func (t *createServiceTask) ID() runner.TaskID {
	return runner.TaskCreateService
}

func (t *createServiceTask) Run(ctx context.Context) (err error) {
	if t.service == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create service %s: %v", t.service.Name, err))
	}()

	v, err := json.Marshal(t.service)
	if err != nil {
		return err
	}

	if err := client.Default().CreateService(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateServiceTask struct {
	service *api.ServiceConfig
}

func UpdateService(service *api.ServiceConfig) runner.Task {
	return &updateServiceTask{
		service: service,
	}
}

func (t *updateServiceTask) ID() runner.TaskID {
	return runner.TaskUpdateService
}

func (t *updateServiceTask) Run(ctx context.Context) (err error) {
	if t.service == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update service %s: %v", t.service.Name, err))
	}()

	v, err := json.Marshal(t.service)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateService(ctx, t.service.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteServiceTask struct {
	service string
}

func DeleteService(service string) runner.Task {
	return &deleteServiceTask{
		service: service,
	}
}

func (t *deleteServiceTask) ID() runner.TaskID {
	return runner.TaskDeleteService
}

func (t *deleteServiceTask) Run(ctx context.Context) (err error) {
	if t.service == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete service %s: %v", t.service, err))
	}()

	if err := client.Default().DeleteService(ctx, t.service); err != nil {
		return err
	}
	return nil
}
