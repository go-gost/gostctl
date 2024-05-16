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

	TaskCreateAdmission TaskID = "api.admission.create"
	TaskUpdateAdmission TaskID = "api.admission.update"
	TaskDeleteAdmission TaskID = "api.admission.delete"

	TaskCreateBypass TaskID = "api.bypass.create"
	TaskUpdateBypass TaskID = "api.bypass.update"
	TaskDeleteBypass TaskID = "api.bypass.delete"

	TaskCreateResolver TaskID = "api.resolver.create"
	TaskUpdateResolver TaskID = "api.resolver.update"
	TaskDeleteResolver TaskID = "api.resolver.delete"

	TaskCreateHosts TaskID = "api.hosts.create"
	TaskUpdateHosts TaskID = "api.hosts.update"
	TaskDeleteHosts TaskID = "api.hosts.delete"

	TaskCreateLimiter TaskID = "api.limiter.create"
	TaskUpdateLimiter TaskID = "api.limiter.update"
	TaskDeleteLimiter TaskID = "api.limiter.delete"
)

type Task interface {
	ID() TaskID
	Run(ctx context.Context) error
}
