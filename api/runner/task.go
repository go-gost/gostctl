package runner

import (
	"context"
)

type TaskID string

const (
	TaskGetConfig  TaskID = "task.api.config.get"
	TaskSaveConfig TaskID = "task.api.config.save"

	TaskCreateService TaskID = "task.api.service.create"
	TaskUpdateService TaskID = "task.api.service.update"
	TaskDeleteService TaskID = "task.api.service.delete"

	TaskCreateChain TaskID = "task.api.chain.create"
	TaskUpdateChain TaskID = "task.api.chain.update"
	TaskDeleteChain TaskID = "task.api.chain.delete"

	TaskCreateHop TaskID = "task.api.hop.create"
	TaskUpdateHop TaskID = "task.api.hop.update"
	TaskDeleteHop TaskID = "task.api.hop.delete"

	TaskCreateAuther TaskID = "task.api.auther.create"
	TaskUpdateAuther TaskID = "task.api.auther.update"
	TaskDeleteAuther TaskID = "task.api.auther.delete"

	TaskCreateAdmission TaskID = "task.api.admission.create"
	TaskUpdateAdmission TaskID = "task.api.admission.update"
	TaskDeleteAdmission TaskID = "task.api.admission.delete"

	TaskCreateBypass TaskID = "task.api.bypass.create"
	TaskUpdateBypass TaskID = "task.api.bypass.update"
	TaskDeleteBypass TaskID = "task.api.bypass.delete"

	TaskCreateResolver TaskID = "task.api.resolver.create"
	TaskUpdateResolver TaskID = "task.api.resolver.update"
	TaskDeleteResolver TaskID = "task.api.resolver.delete"

	TaskCreateHosts TaskID = "task.api.hosts.create"
	TaskUpdateHosts TaskID = "task.api.hosts.update"
	TaskDeleteHosts TaskID = "task.api.hosts.delete"

	TaskCreateLimiter TaskID = "task.api.limiter.create"
	TaskUpdateLimiter TaskID = "task.api.limiter.update"
	TaskDeleteLimiter TaskID = "task.api.limiter.delete"

	TaskCreateObserver TaskID = "task.api.observer.create"
	TaskUpdateObserver TaskID = "task.api.observer.update"
	TaskDeleteObserver TaskID = "task.api.observer.delete"

	TaskCreateRecorder TaskID = "task.api.recorder.create"
	TaskUpdateRecorder TaskID = "task.api.recorder.update"
	TaskDeleteRecorder TaskID = "task.api.recorder.delete"
)

type Task interface {
	ID() TaskID
	Run(ctx context.Context) error
}
