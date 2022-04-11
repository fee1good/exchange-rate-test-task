package jobs

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/fee1good/exchange-rate-test-task/internal/config"
	"github.com/fee1good/exchange-rate-test-task/internal/services/rates"
	"github.com/fee1good/exchange-rate-test-task/internal/services/redis/locker"
)

const uploadRatesJobName = "upload_rates_job"

type uploadRatesJob struct {
	spec                       string
	lockTTL                    time.Duration
	locker                     locker.Locker
	logger                     *zerolog.Logger
	ratesService               *rates.Service
	cryptoSymbols, fiatSymbols []string
}

func newUploadRatesJob(
	cfg *config.Jobs,
	logger *zerolog.Logger,
	locker locker.Locker,
	ratesService *rates.Service,
) job {
	return &uploadRatesJob{
		spec:          cfg.UploadRatesJobSpec,
		lockTTL:       cfg.UploadRatesJobLockTTl,
		logger:        logger,
		locker:        locker,
		ratesService:  ratesService,
		cryptoSymbols: cfg.CryptoSymbols,
		fiatSymbols:   cfg.FiatSymbols,
	}
}

func (j *uploadRatesJob) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobLogger := j.logger.With().Str("job_name", uploadRatesJobName).Logger()
	if _, err := j.locker.Obtain(ctx, uploadRatesJobName, j.lockTTL); err != nil {
		jobLogger.Debug().Err(err).Msg("failed to obtain lock")
		return
	}

	jobLogger.Info().Msg("job start")
	if err := j.ratesService.UploadRates(ctx, j.cryptoSymbols, j.fiatSymbols); err != nil {
		jobLogger.Debug().Err(err).Msg("failed to upload rates")
		return
	}
	jobLogger.Info().Msg("job success")
}

func (j *uploadRatesJob) Spec() string {
	return j.spec
}
