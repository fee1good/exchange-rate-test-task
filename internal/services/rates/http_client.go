package rates

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fee1good/exchange-rate-test-task/internal/entities"
)

const (
	getRatesHTTPPath = "/data/pricemultifull"

	cryptoSymbolsParamKey = "fsyms"
	fiatSymbolsParamKey   = "tsyms"
)

type cryptoCompareAPIClient struct {
	URL        string
	httpClient *http.Client
}

func newClient(url string) *cryptoCompareAPIClient {
	return &cryptoCompareAPIClient{
		URL:        url,
		httpClient: &http.Client{},
	}
}

func (c *cryptoCompareAPIClient) GetPairsRate(ctx context.Context, cryptoSymbols, fiatSymbols []string) (*entities.PairsRate, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.URL, getRatesHTTPPath), nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	query := request.URL.Query()
	query.Add(cryptoSymbolsParamKey, strings.Join(cryptoSymbols, ","))
	query.Add(fiatSymbolsParamKey, strings.Join(fiatSymbols, ","))
	request.URL.RawQuery = query.Encode()

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	var result entities.PairsRate
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	if err := validateRateResponse(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

var MalformedResponseErr = errors.New("response from api is malformed")

func validateRateResponse(res *entities.PairsRate) error {
	if len(res.Raw) == 0 || len(res.Display) == 0 {
		return MalformedResponseErr
	}

	for key := range res.Raw {
		if len(res.Raw[key]) == 0 {
			return MalformedResponseErr
		}

		if displayBucket, ok := res.Display[key]; !ok || len(displayBucket) == 0 {
			return MalformedResponseErr
		}
	}

	return nil
}
