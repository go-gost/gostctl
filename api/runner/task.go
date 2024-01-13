package runner

import (
	"context"
	"time"

	"github.com/go-gost/gui/api"
	"github.com/go-gost/gui/api/client"
)

type Task interface {
	ID() string
	Run(ctx context.Context) error
}

type getConfigTask struct{}

func GetConfigTask() Task {
	return &getConfigTask{}
}

func (t *getConfigTask) ID() string {
	return "get_config"
}

func (t *getConfigTask) Run(ctx context.Context) error {
	cfg, err := client.Default().GetConfig(ctx)
	if err != nil {
		return err
	}

	oldCfg := api.GetConfig()

	for _, service := range cfg.Services {
		if service.Status == nil {
			service.Status = &api.ServiceStatus{}
		}
		if service.Status.Stats == nil {
			continue
		}

		service.Status.Stats.Time = time.Now()

		for _, svc := range oldCfg.Services {
			if svc.Name != service.Name ||
				svc.Status == nil ||
				svc.Status.Stats == nil ||
				svc.Status.CreateTime != service.Status.CreateTime {
				continue
			}

			d := service.Status.Stats.Time.Sub(svc.Status.Stats.Time)
			if d <= 0 {
				continue
			}

			inputRateBytes := int64(service.Status.Stats.InputBytes) - int64(svc.Status.Stats.InputBytes)
			if inputRateBytes < 0 {
				inputRateBytes = 0
			}
			service.Status.Stats.InputRateBytes = uint64(float64(inputRateBytes) / d.Seconds())

			outputRateBytes := int64(service.Status.Stats.OutputBytes) - int64(svc.Status.Stats.OutputBytes)
			if outputRateBytes < 0 {
				outputRateBytes = 0
			}
			service.Status.Stats.OutputRateBytes = uint64(float64(outputRateBytes) / d.Seconds())

			reqRate := int64(service.Status.Stats.TotalConns) - int64(svc.Status.Stats.TotalConns)
			if reqRate < 0 {
				reqRate = 0
			}
			service.Status.Stats.RequestRate = float64(reqRate) / d.Seconds()

			break
		}
	}

	api.SetConfig(cfg)
	return nil
}
