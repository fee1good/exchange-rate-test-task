package jobs

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"

	"github.com/fee1good/exchange-rate-test-task/internal/config"
	"github.com/fee1good/exchange-rate-test-task/internal/services"
	"github.com/fee1good/exchange-rate-test-task/internal/services/redis/locker"
)

type job interface {
	Run()
	Spec() string
}

type Runner interface {
	Start()
	Stop()
}

type jobsRunner struct {
	scheduler *cron.Cron
}

func New(cfg *config.Jobs, logger *zerolog.Logger, locker locker.Locker, services *services.Container) (Runner, error) {
	scheduler := cron.New(cron.WithLogger(newCronLogger(logger)))

	uploadRatesJob := newUploadRatesJob(cfg, logger, locker, services.Rates)
	if _, err := scheduler.AddJob(uploadRatesJob.Spec(), uploadRatesJob); err != nil {
		return nil, fmt.Errorf("failed to add simple job: %w", err)
	}

	return &jobsRunner{scheduler: scheduler}, nil
}

func (j jobsRunner) Start() {
	j.scheduler.Start()
}

func (j jobsRunner) Stop() {
	j.scheduler.Stop()
}
