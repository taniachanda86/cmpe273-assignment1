package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"os"
)

const (
	timeout = time.Duration(time.Second * 10)
)


type RequestMap struct{
	Stocks map[string]float64
}

//Adding this block+++++++++++++++++++
type RequestTradeID struct{
 	TradeId string
 }
////Adding this upper block ++++++++++
type Stock struct {
	List struct {
		Resources []struct {
			Resource struct {
				Fields struct {
					Name    string `json:"name"`
					Price   string `json:"price"`
					Symbol  string `json:"symbol"`
					Ts      string `json:"ts"`
					Type    string `json:"type"`
					UTCTime string `json:"utctime"`
					Volume  string `json:"volume"`
				} `json:"fields"`
			} `json:"resource"`
		} `json:"resources"`
	} `json:"list"`
}

type ResponseList struct{
	TradeID string
	NoOfStocks []int
	Price []float64
	Symbol []string
	Unvested []float64
}

type ResponseTradeID struct{
	Symbol []string
	CurrentPrice []float64
	ChangeInPrice []float64
	Unvested []float64
	NoOfStocks []int
}

//Added+++++++++++++++++++++++++++++++++++++++++++++++++
type Data struct{
				Price float64
				UnvestedAmount float64
				NumberOfStocks int
			}
var data Data
var tradeMap map[string]map[string] Data
var map1 map[string]Data
var tempMap map[string]Data
//Added this upper block++++++++++++++++++++++++++++++++

type Finance int

func (f *Finance) DoTheJob(request *RequestMap, response *ResponseList ) error {

 	///Generating UUID
	tradeID, err := newUUID()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	response.TradeID = string(tradeID)


	//Added this block +++++++++++++++++++++++++++++++++++++++++++
	tradeMap= make(map[string]map[string] Data)	
	tempMap=make(map[string]Data)
	//Added this upper block +++++++++++++++++++++++++++++++++++++


	///Getting the price and number of stocks
	for key, value:= range request.Stocks{
		stockPrice := getQuote(key)
		count:= int(value/stockPrice)
		unvestedAmont := value - float64(count)*stockPrice
		// fmt.Println(int(count), unvestedAmont, stockPrice, value)

		response.Symbol = append(response.Symbol, key)
		response.Price = append(response.Price, stockPrice)
		response.NoOfStocks = append(response.NoOfStocks, int(count))
		response.Unvested = append(response.Unvested, unvestedAmont)

		//Added this block +++++++++++++++++++++++++++++++++++++++++++
		data = Data{stockPrice,unvestedAmont,int(count)}
		// fmt.Println(data)
		tempMap[key]= data
		//Added this upper block +++++++++++++++++++++++++++++++++++++
	}
	tradeMap[tradeID]=tempMap
	return nil
}

func newUUID() (string, error) {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil

}


func getQuote(symbol string) float64 {
	// set http client timeout
	client := http.Client{Timeout: timeout}

	url := fmt.Sprintf("http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json", symbol)
	// fmt.Println(symbol)
	// fmt.Println(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Errorf("Stocks cannot access yahoo finance API: %v", err)
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Stocks cannot read json body: %v", err)
	}
	// fmt.Println(content)
	var stock Stock

	err = json.Unmarshal(content, &stock)
	if err != nil {
		fmt.Errorf("Stocks cannot parse json data: %v", err)
		// return err
	}

	price, err := strconv.ParseFloat(stock.List.Resources[0].Resource.Fields.Price, 64)
	if err != nil {
		fmt.Errorf("Stock price: %v", err)
	}
	// fmt.Println(price)
	return price
}

//Adding this block ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (f *Finance) GetPortfolio(requestId *RequestTradeID, responseId *ResponseTradeID ) error {
	map1 := tradeMap[requestId.TradeId]
	for key,value:= range map1{
		// fmt.Println("The symbl is : ", key)
		updatedPrice := getQuote(key)
		previousPrice := value.Price
		changeInPrice:=updatedPrice - previousPrice

		// fmt.Println("new Price is : ", updatedPrice)
		// fmt.Println("previous Price was : ", previousPrice)
		// fmt.Println("the change is : ", changeInPrice )
		
		responseId.Symbol = append(responseId.Symbol, key)
		responseId.CurrentPrice = append(responseId.CurrentPrice, updatedPrice)
		responseId.ChangeInPrice = append(responseId.ChangeInPrice, changeInPrice)
		responseId.Unvested = append(responseId.Unvested, value.UnvestedAmount)
		responseId.NoOfStocks = append(responseId.NoOfStocks, value.NumberOfStocks)
	}
	return nil

}
//Adding this upper block ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


func main() {

	fin := new(Finance)
	server := rpc.NewServer()
	server.Register(fin)

	l, e := net.Listen("tcp", ":8222")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}