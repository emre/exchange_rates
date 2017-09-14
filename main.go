package main

import (
	"github.com/emre/exchange_rates/open_exchange_rates"
	"fmt"
	"log"
)

func main() {

	c := open_exchange_rates.GetClient("dc5f5be2c4cb44a88d19b6c3bfe05fe5")
	e := GetExchanger(c)

	convertedValue, err := e.Convert(float64(1), "USD", "TRY")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%f", *convertedValue)
}