# cmpe273_assignment1
Connecting to yahoo finance API with Golang.
Instruction to run the Client-Server communication:

Part 1:
  1. First run the Server (go run yahooServer.go).
  2. Run the Client with arguments: Stock symbol and percentage and the total budget allocated (go run yahooClient.go GOOG:60%,YHOO:40% 1000).
  
go run yahooClient.go GOOG:40%,YHOO:60% 1000

    Trade ID:  46c9e018-0ab3-4213-1b6f-7772cc3a1fee
    GOOG:0:$626.909973, YHOO:19:$30.709999, 
    Total unvested amount is:  416.51001900000006
  
Part 2:
  1. Copy the Trade id from the response of Part 1 and pass that as a argument to the client (go run yahooClient.go 46c9e018-0ab3-4213-1b6f-7772cc3a1fee).
  
  go run yahooClient.go 46c9e018-0ab3-4213-1b6f-7772cc3a1fee
    
      GOOG:0:-$626.909973, YHOO:19:-$30.709999, 
      Current market price:  GOOG:$626.909973, YHOO:$30.709999, 
      Total unvested amount is: $ 416.51001900000006
      
  (The "+"/ "-" is not shown in the response of the portfolio because the stock market was closed when I ran it.)
