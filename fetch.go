package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const quotesURLv7 = `https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s`
const quotesURLv7QueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

func main() {
	// fmt.Println("args:", len(os.Args))
	url := fmt.Sprintf(quotesURLv7, "BRK-B,VBK")
	if len(os.Args) > 1 {
		url += quotesURLv7QueryParts
	}
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	parse(body)
}

func parse(body []byte) {
	d := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		panic(err)
	}
	for num, data := range d["quoteResponse"]["result"] {
		fmt.Println("result:", num)
		for k, v := range data {
			fmt.Println(k, ":", v)
		}
	}
	// print_map(d)
}

// func print_map(d map[string]map[string][]map[string]interface{}{}) {
// 	for k, v := range d {
// 		fmt.Println(k)
// 	}
// }
