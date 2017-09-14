package main

import (
	"github.com/emre/exchange_rates/open_exchange_rates"
	"sync"
)

type Exchanger struct {
	Client open_exchange_rates.Client
	Cache InMemoryCache
	mu sync.Mutex
}

func GetExchanger(client open_exchange_rates.Client) Exchanger {

	cache := InMemoryCache{Client: &client}
	e := Exchanger{Client:client, Cache:cache}

	return e
}

func (e *Exchanger) Convert(value float64, baseCurrency string, targetCurrency string) (*float64, error) {
	e.mu.Lock()
	e.Cache.BaseCurrency = baseCurrency
	e.mu.Unlock()
	if baseCurrency == targetCurrency {
		return &value, nil
	}

	rateStore, err := e.Cache.Get(baseCurrency)
	if err != nil {
		return nil, err
	}

	rate := open_exchange_rates.Rate{
		Base:baseCurrency,
		Target:targetCurrency,
		Rate: rateStore.Rates[targetCurrency],
	}

	totalValue := rate.Rate * value
	return &totalValue, nil

}

