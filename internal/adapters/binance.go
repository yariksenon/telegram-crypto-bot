package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"tg-crypto-bot/internal/interfaces"
)

type BinancePriceProvider struct{}

type binanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func NewBinancePriceProvider() interfaces.PriceProvider {
	return &BinancePriceProvider{}
}

func (b *BinancePriceProvider) GetPrice(symbol string) (float64, error) {
	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.New("не удалось получить ответ от Binance")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("не удалось получить ответ от Binance для символа")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("не удалось прочитать тело ответа")
	}

	var data binanceResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, errors.New("не удалось анмаршлить ответ")
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return 0, errors.New("не удалось преобразовать price во float64")
	}

	return price, nil
}
