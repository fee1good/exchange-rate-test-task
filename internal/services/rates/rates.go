package rates

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/multierr"

	"github.com/fee1good/exchange-rate-test-task/internal/config"
	"github.com/fee1good/exchange-rate-test-task/internal/entities"
	"github.com/fee1good/exchange-rate-test-task/internal/repositories/rates"
)

type CryptoCompareAPI interface {
	GetPairsRate(ctx context.Context, cryptoSymbols, fiatSymbols []string) (*entities.PairsRate, error)
}

type RateRepository interface {
	UploadPairs(ctx context.Context, pairs []entities.Pair) error
	GetPairs(ctx context.Context, cryptoSymbols, fiatSymbols []string) ([]entities.Pair, error)
}

type Service struct {
	logger           *zerolog.Logger
	cryptoCompareAPI CryptoCompareAPI
	repository       RateRepository
}

func NewService(cfg *config.Config, logger *zerolog.Logger, repository *rates.Repository) *Service {
	return &Service{
		logger:           logger,
		cryptoCompareAPI: newClient(cfg.CryptoCompareURL),
		repository:       repository,
	}
}

func (s *Service) UploadRates(ctx context.Context, cryptoSymbols, fiatSymbols []string) error {
	result, err := s.cryptoCompareAPI.GetPairsRate(ctx, cryptoSymbols, fiatSymbols)
	if err != nil {
		return err
	}

	return s.repository.UploadPairs(ctx, buildPairs(result))
}

func (s *Service) GetPairsRate(ctx context.Context, cryptoSymbols, fiatSymbols []string) (*entities.PairsRate, error) {
	//todo: first get from db by timestamp
	result, apiErr := s.cryptoCompareAPI.GetPairsRate(ctx, cryptoSymbols, fiatSymbols)
	if apiErr != nil {
		s.logger.Error().Err(apiErr).Msg("failed to get rates from api")
		repoPairs, repoErr := s.repository.GetPairs(ctx, cryptoSymbols, fiatSymbols)
		if repoErr != nil {
			return nil, multierr.Combine(apiErr, repoErr)
		}
		return buildPairsRate(repoPairs), nil
	}

	return result, nil
}

func buildPairsRate(pairs []entities.Pair) *entities.PairsRate {
	result := &entities.PairsRate{
		Raw:     map[string]map[string]entities.RateValues[float64]{},
		Display: map[string]map[string]entities.RateValues[string]{},
	}

	for _, pairRate := range pairs {
		if _, ok := result.Raw[pairRate.CryptoSymbol]; !ok {
			result.Raw[pairRate.CryptoSymbol] = map[string]entities.RateValues[float64]{}
			result.Display[pairRate.CryptoSymbol] = map[string]entities.RateValues[string]{}
		}

		result.Raw[pairRate.CryptoSymbol][pairRate.FiatSymbol] = pairRate.RawRateValues
		result.Display[pairRate.CryptoSymbol][pairRate.FiatSymbol] = pairRate.DisplayRateValues
	}

	return result
}

func buildPairs(resp *entities.PairsRate) []entities.Pair {
	//todo: mocked data, build valid response
	return []entities.Pair{
		{
			CryptoSymbol:      "ETH",
			FiatSymbol:        "USD",
			RawRateValues:     entities.RateValues[float64]{ChangePCT24Hour: 1.1},
			DisplayRateValues: entities.RateValues[string]{ChangePCT24Hour: "1.1"},
		},
	}
}
