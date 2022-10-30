package main

// usage:
//   fetch --dur <delay> [--fake]
//
import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	etfs "tgphelps.com/etfs/pkg"
)

type stock_data struct {
	market     string
	symbol     string
	price      float64
	ave_50day  float64
	ave_200day float64
	prev_close float64
	open       float64
	low        float64
	high       float64
	low_52w    float64
	high_52w   float64
}

// Constant values

const HDG1 = " sym  price  ave 50 ave200 close   open   low    high  low52w hi52w   chg"
const HDG2 = "===== ====== ====== ====== ====== ====== ====== ====== ====== ====== ====="
const URL1 = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s"
const URL2 = "&range=1d&interval=5m&indicators=close&includeTimestamps=false" +
	"&includePrePost=false&corsDomain=finance.yahoo.com" +
	"&.tsrc=finance"

// Global variables

var g_fake_data bool
var g_interval time.Duration

func init() {
	flag.DurationVar(&g_interval, "dur", 5*time.Minute, "web fetch delay")
	flag.BoolVar(&g_fake_data, "fake", false, "use fake data")
}

func main() {
	var body []byte
	var err error

	flag.Parse()
	fmt.Println("Interval:", g_interval)
	fmt.Println("Fake?:", g_fake_data)
	fmt.Println("flag args:", flag.Arg(0))

	// This loops until the user interrupts it with ctrl+C.
	for {
		fmt.Println("")
		dt := time.Now()
		fmt.Println("Time:", dt.Format("15:04:05"))
		fmt.Println(HDG1)
		fmt.Println(HDG2)
		if g_fake_data {
			body, err = os.ReadFile("body.txt")
			check(err)
		} else {
			body = etfs.Fetch_stock_data(flag.Arg(0))
		}
		if parse_and_print(body) {
			break
		}
		time.Sleep(g_interval)
	}
}

func parse_and_print(body []byte) bool {
	var market_state string
	d := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		fmt.Println("json.Unmarshal error:")
		fmt.Println(err)
		fmt.Println(body)
		panic("STOP")
	}
	for _, dict_data := range d["quoteResponse"]["result"] {
		var data stock_data
		// fmt.Println("result:", num)
		build_data(dict_data, &data)
		print_stock_data(&data)
		market_state = data.market
	}
	fmt.Println("Market is ", market_state)
	if market_state == "POST" {
		return true
	} else {
		return false
	}
}

func build_data(dict_data map[string]interface{}, data *stock_data) {
	data.market = dict_data["marketState"].(string)
	data.symbol = dict_data["symbol"].(string)
	data.price = dict_data["regularMarketPrice"].(float64)
	data.ave_50day = dict_data["fiftyDayAverage"].(float64)
	data.ave_200day = dict_data["twoHundredDayAverage"].(float64)
	data.prev_close = dict_data["regularMarketPreviousClose"].(float64)
	data.open = dict_data["regularMarketOpen"].(float64)
	data.low = dict_data["regularMarketDayLow"].(float64)
	data.high = dict_data["regularMarketDayHigh"].(float64)
	data.low_52w = dict_data["fiftyTwoWeekLow"].(float64)
	data.high_52w = dict_data["fiftyTwoWeekHigh"].(float64)
}

func print_stock_data(data *stock_data) {
	fmt.Printf("%-5s %6.2f %6.2f %6.2f %6.2f %6.2f %6.2f %6.2f %6.2f %6.2f %5.2f\n",
		data.symbol, data.price, data.ave_50day, data.ave_200day,
		data.prev_close, data.open, data.low, data.high,
		data.low_52w, data.high_52w, data.price-data.prev_close)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
