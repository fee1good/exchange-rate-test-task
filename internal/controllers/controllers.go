package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/fee1good/exchange-rate-test-task/internal/services/rates"
)

type Container struct {
	logger       *zerolog.Logger
	ratesService *rates.Service
}

func NewHTTPContainer(logger *zerolog.Logger, ratesService *rates.Service) *Container {
	return &Container{
		logger:       logger,
		ratesService: ratesService,
	}
}

func (c *Container) Mux() *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	mux.Get("/service/price", c.getRate)
	return mux
}

func (c *Container) getRate(w http.ResponseWriter, r *http.Request) {
	cryptoSymbols := strings.Split(chi.URLParam(r, "tsyms"), ",")
	fiatSymbols := strings.Split(chi.URLParam(r, "fsyms"), ",")
	if len(cryptoSymbols) == 0 || len(fiatSymbols) == 0 {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return
	}

	pairsRate, err := c.ratesService.GetPairsRate(r.Context(), cryptoSymbols, fiatSymbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(pairsRate)
	if err != nil {
		http.Error(w, "failed to serialize response", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(responseBytes); err != nil {
		c.logger.Error().Err(err).Msg("failed to write response")
	}
}
