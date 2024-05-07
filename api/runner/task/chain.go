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

type createChainTask struct {
	chain *api.ChainConfig
}

func CreateChain(chain *api.ChainConfig) runner.Task {
	return &createChainTask{
		chain: chain,
	}
}

func (t *createChainTask) ID() runner.TaskID {
	return runner.TaskCreateChain
}

func (t *createChainTask) Run(ctx context.Context) (err error) {
	if t.chain == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("create chain %s: %v", t.chain.Name, err))
	}()

	v, err := json.Marshal(t.chain)
	if err != nil {
		return err
	}

	if err := client.Default().CreateChain(ctx, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type updateChainTask struct {
	chain *api.ChainConfig
}

func UpdateChain(chain *api.ChainConfig) runner.Task {
	return &updateChainTask{
		chain: chain,
	}
}

func (t *updateChainTask) ID() runner.TaskID {
	return runner.TaskUpdateChain
}

func (t *updateChainTask) Run(ctx context.Context) (err error) {
	if t.chain == nil {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("update chain %s: %v", t.chain.Name, err))
	}()

	v, err := json.Marshal(t.chain)
	if err != nil {
		return err
	}

	if err := client.Default().UpdateChain(ctx, t.chain.Name, bytes.NewReader(v)); err != nil {
		return err
	}
	return nil
}

type deleteChainTask struct {
	chain string
}

func DeleteChain(chain string) runner.Task {
	return &deleteChainTask{
		chain: chain,
	}
}

func (t *deleteChainTask) ID() runner.TaskID {
	return runner.TaskDeleteChain
}

func (t *deleteChainTask) Run(ctx context.Context) (err error) {
	if t.chain == "" {
		return nil
	}

	defer func() {
		slog.With("kind", "task", "task", t.ID()).DebugContext(ctx, fmt.Sprintf("delete chain %s: %v", t.chain, err))
	}()

	if err := client.Default().DeleteChain(ctx, t.chain); err != nil {
		return err
	}
	return nil
}
