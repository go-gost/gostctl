package runner

import (
	"context"
)

type TaskID string

const (
	TaskGetConfig  TaskID = "api.config.get"
	TaskSaveConfig TaskID = "api.config.save"

	TaskCreateService TaskID = "api.service.create"
	TaskUpdateService TaskID = "api.service.update"
	TaskDeleteService TaskID = "api.service.delete"
)

type Task interface {
	ID() TaskID
	Run(ctx context.Context) error
}
