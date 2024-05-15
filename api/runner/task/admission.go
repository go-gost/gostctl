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

type createAdmissionTask struct {
	admission *api.AdmissionConfig
}

func CreateAdmission(admission *api.AdmissionConfig) runner.Task {
	return &createAdmissionTask{
		admission: admission,
	}
}

func (t *createAdmissionTask) ID() runner.TaskID {
	return runner.TaskCreateAdmission
}

func (t *createAdmissionTask) Run(ctx context.Context) (err error) {
	if t.admission == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create admission %s: %v", t.admission.Name, err))
	}()

	v, err := json.Marshal(t.admission)
	if err != nil {
		return err
	}

	if err := client.Default().CreateAdmission(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateAdmissionTask struct {
	admission *api.AdmissionConfig
}

func UpdateAdmission(admission *api.AdmissionConfig) runner.Task {
	return &updateAdmissionTask{
		admission: admission,
	}
}

func (t *updateAdmissionTask) ID() runner.TaskID {
	return runner.TaskUpdateAdmission
}

func (t *updateAdmissionTask) Run(ctx context.Context) (err error) {
	if t.admission == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update admission %s: %v", t.admission.Name, err))
	}()

	v, err := json.Marshal(t.admission)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateAdmission(ctx, t.admission.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteAdmissionTask struct {
	admission string
}

func DeleteAdmission(admission string) runner.Task {
	return &deleteAdmissionTask{
		admission: admission,
	}
}

func (t *deleteAdmissionTask) ID() runner.TaskID {
	return runner.TaskDeleteAdmission
}

func (t *deleteAdmissionTask) Run(ctx context.Context) (err error) {
	if t.admission == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete admission %s: %v", t.admission, err))
	}()

	if err := client.Default().DeleteAdmission(ctx, t.admission); err != nil {
		return err
	}
	return nil
}
