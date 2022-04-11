package config

import (
	"fmt"
	"time"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Port             string `envconfig:"PORT"`
	LogLevel         string `envconfig:"LOG_LEVEL"`
	PGDSN            string `envconfig:"PG_DSN"`
	CryptoCompareURL string `envconfig:"CRYPTO_COMPARE_API_URL"`

	*Jobs
	*Redis
}

type Jobs struct {
	UploadRatesJobSpec    string        `envconfig:"UPLOAD_RATES_JOB_SPEC"`
	UploadRatesJobLockTTl time.Duration `envconfig:"UPLOAD_RATES_JOB_LOCK_TTL,default=1m"`
	CryptoSymbols         []string      `envconfig:"CRYPTO_SYMBOLS"`
	FiatSymbols           []string      `envconfig:"FIAT_SYMBOLS"`
}

type Redis struct {
	URL        string `envconfig:"REDIS_URL"`
	MasterName string `envconfig:"REDIS_MASTER_NAME"`
	Password   string `envconfig:"REDIS_PASSWORD,optional"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Init(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
