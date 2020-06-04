package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	bitco "github.com/hypoballad/bitcocheck"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:50051", "server address")
var actualMode = flag.Bool("actual", false, "actual mode")
var debugMode = flag.Bool("debug", false, "mode debug")

func BuyRateBtc(conn *grpc.ClientConn, value string) (*bitco.ExchangeOrdersRateItem, error) {
	return OrderRate(conn, bitco.Buy, bitco.Price, value)
}

func SellRateBtc(conn *grpc.ClientConn, value string) (*bitco.ExchangeOrdersRateItem, error) {
	return OrderRate(conn, bitco.Sell, bitco.Amount, value)
}

func OrderRate(conn *grpc.ClientConn, order bitco.OrderType, amountprice bitco.AmountPriceType, value string) (*bitco.ExchangeOrdersRateItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var item *bitco.ExchangeOrdersRateItem
	var param bitco.ExchangeOrdersRateParam
	param.OrderType = order.String()
	param.Pair = bitco.Btcjpy.String()
	param.Amountprice = amountprice.String()
	param.Value = value
	item, err := c.ExchangeOrdersRate(ctx, &param)
	if err != nil {
		return item, err
	}
	return item, nil
}

func SalesRate(conn *grpc.ClientConn) (*bitco.RatePairItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var item *bitco.RatePairItem
	var param bitco.RatePairParams
	param.Pair = bitco.Btcjpy.String()
	item, err := c.RatePair(ctx, &param)
	if err != nil {
		return item, err
	}
	return item, nil
}

func AccountsBalance(conn *grpc.ClientConn) (*bitco.AccountsBalanceItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var item *bitco.AccountsBalanceItem
	var in bitco.Empty
	item, err := c.AccountsBalance(ctx, &in)
	if err != nil {
		return item, err
	}
	return item, nil
}

func Accounts(conn *grpc.ClientConn) (*bitco.AccountsItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var item *bitco.AccountsItem
	var in bitco.Empty
	item, err := c.Accounts(ctx, &in)
	if err != nil {
		return item, err
	}
	return item, nil
}

func Trades(conn *grpc.ClientConn) (*bitco.TradesItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	in := &bitco.TradesParams{Pair: "btc_jpy"}
	item, err := c.Trades(ctx, in)
	if err != nil {
		return item, err
	}
	return item, nil
}

func debugJson(v interface{}) {
	b, err := json.MarshalIndent(v, "", "	")
	if err != nil {
		fmt.Println("debug print error: ", err)
	}
	fmt.Println(string(b))
}
func job(addr string, debug bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// account balance
	balance, err := AccountsBalance(conn)
	if err != nil {
		log.Println("accounts balance error: ", err)
	}
	btc, err := strconv.ParseFloat(balance.Btc, 32)
	if err != nil {
		log.Println("btc convert error:", err)
	}
	fmt.Printf("myaccount yen: %s, btc: %f\n", balance.Jpy, btc)

	salesrate, err := SalesRate(conn)
	if err != nil {
		log.Println("sales rate error:", err)
	}

	// fmt.Println(salesrate.Rate)
	buyrate, err := BuyRateBtc(conn, balance.Jpy)
	if err != nil {
		log.Println("buy rate error:", err)
	}
	// debugJson(item)
	fmt.Printf("1btc=%s, buy:  1btc=%s, %s[yen]=%s[btc]\n", salesrate.Rate, buyrate.Rate, buyrate.Price, buyrate.Amount)

	sellrate, err := SellRateBtc(conn, balance.Btc)
	if err != nil {
		log.Println("buy rate error:", err)
	}
	// debugJson(sellrate)
	fmt.Printf("1btc=%s, sell: 1btc=%s, %s[yen]=%s[btc]\n", salesrate.Rate, sellrate.Rate, sellrate.Price, sellrate.Amount)
}

func main() {
	flag.Parse()
	// conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// fmt.Println("buy rate")
	// item, err := OrderRate(conn, bitco.Buy)
	// if err != nil {
	// 	log.Fatalln("buy rate error: ", err)
	// }
	// debugJson(item)
	// fmt.Println("sell rate")
	// item, err = OrderRate(conn, bitco.Sell)
	// if err != nil {
	// 	log.Fatalln("sell rate error: ", err)
	// }
	// debugJson(item)
	// fmt.Println("account balance")
	// accBalance, err := AccountsBalance(conn)
	// if err != nil {
	// 	log.Fatalln("accounts balance error: ", err)
	// }
	// debugJson(accBalance)
	// fmt.Println("accounts")
	// acc, err := Accounts(conn)
	// if err != nil {
	// 	log.Fatalln("accounts error: ", err)
	// }
	// debugJson(acc)

	// fmt.Println("trades")
	// trade, err := Trades(conn)
	// if err != nil {
	// 	log.Fatalln("trades error: ", err)
	// }
	// debugJson(trade)
	job(*addr, *debugMode)
	os.Exit(0)
}
