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

	TaskCreateChain TaskID = "api.chain.create"
	TaskUpdateChain TaskID = "api.chain.update"
	TaskDeleteChain TaskID = "api.chain.delete"

	TaskCreateHop TaskID = "api.hop.create"
	TaskUpdateHop TaskID = "api.hop.update"
	TaskDeleteHop TaskID = "api.hop.delete"

	TaskCreateAuther TaskID = "api.auther.create"
	TaskUpdateAuther TaskID = "api.auther.update"
	TaskDeleteAuther TaskID = "api.auther.delete"
)

type Task interface {
	ID() TaskID
	Run(ctx context.Context) error
}
