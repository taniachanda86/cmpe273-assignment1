package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"net"
	"net/rpc/jsonrpc"
)

type RequestMap struct{
	Stocks map[string]float64
}

type ResponseList struct{
	TradeID string
	NoOfStocks []int
	Price []float64
	Symbol []string
	Unvested []float64
}

type RequestTradeID struct{
 	TradeId string
 }
type ResponseTradeID struct{
	Symbol []string
	CurrentPrice []float64
	ChangeInPrice []float64
	Unvested []float64
	NoOfStocks []int
}



func main() {

	stockMap := make(map[string]float64)

	conn, err := net.Dial("tcp", "localhost:8222")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := jsonrpc.NewClient(conn)
// fmt.Println(os.Args)
if len(os.Args)>2{

	paraOne := strings.Split(os.Args[1], ",")
	paraTwo := os.Args[2]
	budget, err := strconv.Atoi(paraTwo)
	if (err!=nil){
		panic("Fatal Error!")
	}
	for i:=0; i< len(paraOne); i++{
	
		eachSym := strings.Split(paraOne[i],":")		
		percent:= strings.TrimSuffix(eachSym[1],"%")
		percentage, err1 := strconv.Atoi(percent)
		
		if (err1!=nil){
			panic("Fatal Error!")
		}
		allotedMoney := float64(budget * percentage/100)
		// request = &RequestMap{eachSym[0]:allotedMoney}
		stockMap[eachSym[0]] = allotedMoney
	}
	var request *RequestMap
	var response ResponseList
	request = &RequestMap{stockMap}
    
    err = c.Call("Finance.DoTheJob", request, &response)
	if err != nil {
		fmt.Errorf("finance error:", err)
	}

	// SaveData[response.TradeID] = [response.Symbol[], strconv.Itoa(response.NoOfStocks[]),]


	fmt.Println("Trade ID: ", response.TradeID)
	var finalString string
	var totalUnvested float64
	for i:=range response.Symbol{

		finalString += fmt.Sprintf("%s:%d:$%g, ", response.Symbol[i],response.NoOfStocks[i],response.Price[i])
		totalUnvested +=response.Unvested[i]
	}
	fmt.Println(finalString)
	fmt.Println("Total unvested amount is: ", totalUnvested)

///Adding this block++++++++++++++++++++++++++++++++++++++++++++++++
}else{
	// fmt.Println(len(os.Args))
	getTradeID:=os.Args[1]
	// fmt.Println(getTradeID)
	var requestId *RequestTradeID
	requestId = &RequestTradeID{getTradeID}
	var responseId ResponseTradeID
	err = c.Call("Finance.GetPortfolio", requestId, &responseId)
	if err != nil {
		fmt.Errorf("finance error:", err)
	}
	// fmt.Println("%v\n", responseId)
	var finalPortfolioString, finalCurrentValue string
	var totalPortfolioUnvested float64
	for i:=range responseId.Symbol{
		if responseId.ChangeInPrice[i]>0{

			finalPortfolioString += fmt.Sprintf("%s:%d:+$%g, ", responseId.Symbol[i],responseId.NoOfStocks[i],responseId.CurrentPrice[i])
		}else{
			finalPortfolioString += fmt.Sprintf("%s:%d:-$%g, ", responseId.Symbol[i],responseId.NoOfStocks[i],responseId.CurrentPrice[i])
		}

		finalCurrentValue += fmt.Sprintf("%s:$%g, ",responseId.Symbol[i], responseId.CurrentPrice[i])
		totalPortfolioUnvested +=responseId.Unvested[i]
	}
	fmt.Println(finalPortfolioString)
	fmt.Println("Current market price: ", finalCurrentValue)
	fmt.Println("Total unvested amount is: $", totalPortfolioUnvested)


}
////Adding the upper block +++++++++++++++++++++++++++++++++++++++++++
	
}
