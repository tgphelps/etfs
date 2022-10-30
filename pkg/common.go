package etfs

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type holding struct {
	sym        string
	shares     int
	total_cost float32
	ave_cost   float32
}

func (h *holding) buy_shares(shares int, price float32) {
	h.shares += shares
	h.total_cost += float32(shares) * price
	h.ave_cost = h.total_cost / float32(h.shares)
}

func (h *holding) sell_shares(shares int, price float32) {
	if shares > h.shares {
		log.Panic("selling more shares than we own")
	}
	h.shares -= shares
	h.total_cost = float32(h.shares) * h.ave_cost
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
