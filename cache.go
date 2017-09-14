package main

import (
	"github.com/emre/exchange_rates/open_exchange_rates"
	"log"
	"time"
	"sync"
)

var TTL = int32(60 * 60)
var cache = map[string]*open_exchange_rates.RateStore{}


type InMemoryCache struct {
	RateStore *open_exchange_rates.RateStore
	Client *open_exchange_rates.Client
	BaseCurrency string
	mu sync.Mutex
}

func (i *InMemoryCache) Get(baseCurrency string) (*open_exchange_rates.RateStore, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	cachedData := cache[baseCurrency]

	if cachedData == nil {
		log.Printf("Cache miss: %s", baseCurrency)
		err := i.Save(baseCurrency)
		if err != nil {
			return nil, err
		}
		return cache[baseCurrency], nil
	}

	currentTimestamp := int32(time.Now().Unix())
	if cachedData.Timestamp + TTL < currentTimestamp {
		log.Printf("%s cache is expired. Revalidating.", baseCurrency)
		err := i.Save(baseCurrency)
		if err != nil {
			return nil, err
		}
		return cache[baseCurrency], nil
	}
	log.Printf("Cache hit: %s", baseCurrency)

	return cachedData, nil
}

func (i *InMemoryCache) Save(baseCurrency string) error {
	if cache == nil {
		cache = map[string]*open_exchange_rates.RateStore{}
	}
	rateStore, err := i.GetFreshData(baseCurrency)
	if err != nil {
		return err
	}
	cache[baseCurrency] = rateStore
	return nil
}


func (i *InMemoryCache) GetFreshData(baseCurrency string) (rateStore *open_exchange_rates.RateStore, err error) {
	c := open_exchange_rates.GetClient(i.Client.AppId)
	rateStore, err = c.GetRates(baseCurrency)
	return rateStore, err

}