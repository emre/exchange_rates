package open_exchange_rates

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"errors"
)

const (
	apiBaseUrl = "https://openexchangerates.org/api/latest.json"
)


type Client struct {
	HttpClient *http.Client
	AppId string
	BaseUrl *url.URL

}


func (c *Client) GetRates(baseCurrency string)  (rates *RateStore, err error) {
	urlPath := fmt.Sprintf("%s?app_id=%s&base=%s", c.BaseUrl, c.AppId, baseCurrency)
	url, err := url.Parse(urlPath)

	if err != nil {
		return nil, err
	}

	url = c.BaseUrl.ResolveReference(url)
	log.Println("Sending request to:", url)

	response, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if c := response.StatusCode; c >= 400 && c <= 500 {
		errorResponse := ErrorResponse{}
		json.Unmarshal(body, &errorResponse)
		return nil, errors.New(errorResponse.Error())
	}
	rateStore := RateStore{}
	err = json.Unmarshal(body, &rateStore)

	return &rateStore, nil

}

func (c *Client) GetRate(baseCurrency string, targetCurrency string) Rate {
	rateStore, err := c.GetRates(baseCurrency)
	if err != nil {
		log.Fatal(err)
	}
	return Rate{
		Base:baseCurrency,
		Target:targetCurrency,
		Rate: rateStore.Rates[targetCurrency],
	}
}

type ApiResponse  struct {
	TimeStamp int64
	Base string
	Rates map[string]float64
}

func GetClient(appId string) (Client) {

	httpClient := http.DefaultClient
	baseURL, err := url.Parse(apiBaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	c := Client{
		HttpClient: httpClient,
		AppId: appId,
		BaseUrl: baseURL,

	}

	return c
}