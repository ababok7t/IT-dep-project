package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type Coin struct {
	Id               string `json:"id"`
	Symbol           string `json:"symbol"`
	Name             string `json:"name"`
	PriceUsd         string `json:"price_usd"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange7D  string `json:"percent_change_7d"`
}

func getCoinInfo(id string) (Coin, error) {
	var coin []Coin

	response, responseError := http.Get(fmt.Sprint("https://api.coinlore.net/api/ticker/?id=", id))
	if responseError != nil {
		return Coin{}, responseError
	}
	defer response.Body.Close()

	byteSlice := make([]byte, 1024)
	n, _ := response.Body.Read(byteSlice)
	json.Unmarshal(byteSlice[:n], &coin)

	return coin[0], nil
}

func getCoinForecast(id string) (string, error) {
	coin, gettingError := getCoinInfo(id)
	if gettingError != nil {
		return "", gettingError
	}

	priceNow, _ := strconv.ParseFloat(coin.PriceUsd, 64)
	change1H, _ := strconv.ParseFloat(coin.PercentChange1H, 64)
	change24H, _ := strconv.ParseFloat(coin.PercentChange24H, 64)
	change7D, _ := strconv.ParseFloat(coin.PercentChange7D, 64)

	price1H := priceNow / (1 + change1H/100)
	price24H := priceNow / (1 + change24H/100)
	price7D := priceNow / (1 + change7D/100)

	logProfit1 := math.Log(priceNow / price1H)
	logProfit2 := math.Log(price24H / price1H)
	logProfit3 := math.Log(price7D / price24H)

	midProfit := (logProfit1 + logProfit2 + logProfit3) / 3

	vol1 := math.Pow(logProfit1-midProfit, 2)
	vol2 := math.Pow(logProfit2-midProfit, 2)
	vol3 := math.Pow(logProfit3-midProfit, 2)

	volatility := math.Sqrt((vol1 + vol2 + vol3) / 3)

	coinForecast := priceNow*math.Exp(midProfit-1/2*math.Pow(volatility, 2)) + volatility*0.0000001
	return fmt.Sprintf("%.10f", coinForecast), nil
}

func main() {
	btc, err1 := getCoinInfo("90")
	fmt.Println(btc, err1)
	forecast, err2 := getCoinForecast("90")
	fmt.Println(forecast, err2)
}
