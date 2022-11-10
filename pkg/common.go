package etfs

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Holding struct {
	Sym        string
	Shares     int
	Total_cost float64
	Ave_cost   float64
}

func (h *Holding) Buy_shares(shares int, price float64) {
	h.Shares += shares
	h.Total_cost += float64(shares) * price
	h.Ave_cost = h.Total_cost / float64(h.Shares)
}

func (h *Holding) Sell_shares(shares int, price float64) float64 {
	if shares > h.Shares {
		log.Panic("selling more shares than we own")
	}
	h.Shares -= shares
	h.Total_cost = float64(h.Shares) * h.Ave_cost
	return float64(shares) * (price - h.Ave_cost)
}

const URL1 = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s"
const URL2 = "&range=1d&interval=5m&indicators=close&includeTimestamps=false" +
	"&includePrePost=false&corsDomain=finance.yahoo.com" +
	"&.tsrc=finance"

func Fetch_stock_data(stock_list string) []byte {
	// fmt.Println("stock_list: ", stock_list)
	url := fmt.Sprintf(URL1, stock_list)
	url += URL2
	response, err := http.Get(url)
	check(err)
	body, err := io.ReadAll(response.Body)
	check(err)
	response.Body.Close()
	return body
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
