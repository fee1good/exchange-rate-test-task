package rates

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/fee1good/exchange-rate-test-task/internal/entities"
	ratestestmocks "github.com/fee1good/exchange-rate-test-task/test/mocks/packages/ratesservice"
)

func TestService_GetPairsRate(t *testing.T) {
	type mocks struct {
		cryptoCompareAPI *ratestestmocks.MockCryptoCompareAPI
		rateRepository   *ratestestmocks.MockRateRepository
	}

	testRawData := map[string]map[string]entities.RateValues[float64]{
		"BTC": {
			"USD": entities.RateValues[float64]{
				Change24Hour:    2.1,
				ChangePCT24Hour: 2.1,
				Open24Hour:      2.1,
				Volume24Hour:    2.1,
				Volume24HourTo:  2.1,
				Low24Hour:       2.1,
				High24Hour:      2.1,
				Price:           2.1,
				Supply:          2.1,
				MKTCap:          2.1,
			},
			"EUR": entities.RateValues[float64]{
				Change24Hour:    2.2,
				ChangePCT24Hour: 2.2,
				Open24Hour:      2.2,
				Volume24Hour:    2.2,
				Volume24HourTo:  2.2,
				Low24Hour:       2.2,
				High24Hour:      2.2,
				Price:           2.2,
				Supply:          2.2,
				MKTCap:          2.2,
			},
		},
		"ETH": {
			"USD": entities.RateValues[float64]{
				Change24Hour:    1.1,
				ChangePCT24Hour: 1.1,
				Open24Hour:      1.1,
				Volume24Hour:    1.1,
				Volume24HourTo:  1.1,
				Low24Hour:       1.1,
				High24Hour:      1.1,
				Price:           1.1,
				Supply:          1.1,
				MKTCap:          1.1,
			},
			"EUR": entities.RateValues[float64]{
				Change24Hour:    1.2,
				ChangePCT24Hour: 1.2,
				Open24Hour:      1.2,
				Volume24Hour:    1.2,
				Volume24HourTo:  1.2,
				Low24Hour:       1.2,
				High24Hour:      1.2,
				Price:           1.2,
				Supply:          1.2,
				MKTCap:          1.2,
			},
		},
	}

	testDisplayData := map[string]map[string]entities.RateValues[string]{
		"BTC": {
			"USD": entities.RateValues[string]{
				Change24Hour:    "2.1",
				ChangePCT24Hour: "2.1",
				Open24Hour:      "2.1",
				Volume24Hour:    "2.1",
				Volume24HourTo:  "2.1",
				Low24Hour:       "2.1",
				High24Hour:      "2.1",
				Price:           "2.1",
				Supply:          "2.1",
				MKTCap:          "2.1",
			},
			"EUR": entities.RateValues[string]{
				Change24Hour:    "2.2",
				ChangePCT24Hour: "2.2",
				Open24Hour:      "2.2",
				Volume24Hour:    "2.2",
				Volume24HourTo:  "2.2",
				Low24Hour:       "2.2",
				High24Hour:      "2.2",
				Price:           "2.2",
				Supply:          "2.2",
				MKTCap:          "2.2",
			},
		},
		"ETH": {
			"USD": entities.RateValues[string]{
				Change24Hour:    "1.1",
				ChangePCT24Hour: "1.1",
				Open24Hour:      "1.1",
				Volume24Hour:    "1.1",
				Volume24HourTo:  "1.1",
				Low24Hour:       "1.1",
				High24Hour:      "1.1",
				Price:           "1.1",
				Supply:          "1.1",
				MKTCap:          "1.1",
			},
			"EUR": entities.RateValues[string]{
				Change24Hour:    "1.2",
				ChangePCT24Hour: "1.2",
				Open24Hour:      "1.2",
				Volume24Hour:    "1.2",
				Volume24HourTo:  "1.2",
				Low24Hour:       "1.2",
				High24Hour:      "1.2",
				Price:           "1.2",
				Supply:          "1.2",
				MKTCap:          "1.2",
			},
		},
	}

	cases := []struct {
		name                       string
		cryptoSymbols, fiatSymbols []string
		result                     *entities.PairsRate
		isError                    bool
		errMessage                 string
		mocks                      func(ctx context.Context, mocks mocks)
	}{
		{
			name:          "OK with api return data",
			cryptoSymbols: []string{"BTC", "ETH"},
			fiatSymbols:   []string{"USD", "EUR"},
			result: &entities.PairsRate{
				Raw:     testRawData,
				Display: testDisplayData,
			},
			mocks: func(ctx context.Context, mocks mocks) {
				mocks.cryptoCompareAPI.
					EXPECT().
					GetPairsRate(ctx, []string{"BTC", "ETH"}, []string{"USD", "EUR"}).
					Return(&entities.PairsRate{
						Raw:     testRawData,
						Display: testDisplayData,
					}, nil).
					Times(1)
			},
		},
		{
			name:          "OK with api don't return data but repository return data",
			cryptoSymbols: []string{"BTC", "ETH"},
			fiatSymbols:   []string{"USD", "EUR"},
			result: &entities.PairsRate{
				Raw:     testRawData,
				Display: testDisplayData,
			},
			mocks: func(ctx context.Context, mocks mocks) {
				mocks.cryptoCompareAPI.
					EXPECT().
					GetPairsRate(ctx, []string{"BTC", "ETH"}, []string{"USD", "EUR"}).
					Return(nil, errors.New("custom error")).
					Times(1)
				mocks.rateRepository.
					EXPECT().
					GetPairs(ctx, []string{"BTC", "ETH"}, []string{"USD", "EUR"}).
					Return([]entities.Pair{
						{"ETH", "USD", testRawData["ETH"]["USD"], testDisplayData["ETH"]["USD"]},
						{"ETH", "EUR", testRawData["ETH"]["EUR"], testDisplayData["ETH"]["EUR"]},
						{"BTC", "USD", testRawData["BTC"]["USD"], testDisplayData["BTC"]["USD"]},
						{"BTC", "EUR", testRawData["BTC"]["EUR"], testDisplayData["BTC"]["EUR"]},
					}, nil).
					Times(1)
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				ctx      = context.Background()
				logger   = zerolog.Nop()
				mockCtrl = gomock.NewController(t)
			)
			defer mockCtrl.Finish()

			mocks := mocks{
				cryptoCompareAPI: ratestestmocks.NewMockCryptoCompareAPI(mockCtrl),
				rateRepository:   ratestestmocks.NewMockRateRepository(mockCtrl),
			}
			testCase.mocks(ctx, mocks)

			service := &Service{
				logger:           &logger,
				repository:       mocks.rateRepository,
				cryptoCompareAPI: mocks.cryptoCompareAPI,
			}

			pairsRate, err := service.GetPairsRate(ctx, testCase.cryptoSymbols, testCase.fiatSymbols)
			if testCase.isError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.result, pairsRate)
		})
	}
}
