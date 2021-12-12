package main

// fetch --dur 5m --fake
// init - build flag data
//
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type stock_data struct {
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
const FAKE_DATA = true

const HDG1 = " sym  price  ave 50 ave200 close   open   low    high  low52w hi52w"
const HDG2 = "===== ====== ====== ====== ====== ====== ====== ====== ====== ======"
const URL1 = `https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s`
const URL2 = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

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

	for {
		fmt.Println(HDG1)
		fmt.Println(HDG2)
		if FAKE_DATA {
			fmt.Println("FAKE data")
			body, err = os.ReadFile("body.txt")
			check(err)
		} else {
			url := fmt.Sprintf(URL1, "BRK-B,VBK")
			url += URL2
			response, err := http.Get(url)
			check(err)
			body, err = ioutil.ReadAll(response.Body)
			check(err)
			response.Body.Close()
		}
		// fmt.Println(string(body))
		// err = os.WriteFile("body.txt", body, 0644)
		// check(err)
		parse_and_print(body)
		time.Sleep(g_interval)
	}
}

func parse_and_print(body []byte) {
	d := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal(body, &d)
	check(err)
	for num, dict_data := range d["quoteResponse"]["result"] {
		var data stock_data
		fmt.Println("result:", num)
		build_data(dict_data, &data)
		print_stock_data(&data)
		// keys := make([]string, 0, len(data))
		// for k := range data {
		// 	keys = append(keys, k)
		// }
		// sort.Strings(keys)
		// for _, k := range keys {
		// 	fmt.Printf("%s : %T : ", k, data[k])
		// 	fmt.Println(data[k])
		// }
	}
	// print_map(d)
}

func build_data(dict_data map[string]interface{}, data *stock_data) {
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
	fmt.Println(data.symbol, data.price, data.ave_50day, data.ave_200day,
		data.prev_close, data.open, data.low, data.high,
		data.low_52w, data.high_52w)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
