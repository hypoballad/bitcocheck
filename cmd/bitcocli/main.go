package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	bitco "github.com/hypoballad/bitcocheck"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:50051", "server address")
var modeDebug = flag.Bool("debug", false, "debug mode")

// func CheckAll(c bitco.CoincheckClient, ctx context.Context) error {
// 	log.Println("-- coin check ticker --")
// 	in := &bitco.Empty{}
// 	item, err := c.Ticker(ctx, in)
// 	if err != nil {
// 		return err
// 	}
// 	b, err := json.Marshal(item)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(b)
// 	fmt.Println()

// 	return nil
// }

func ticker(conn *grpc.ClientConn) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("-- coin check ticker --")
	in := &bitco.Empty{}
	item, err := c.Ticker(ctx, in)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(item)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	fmt.Println()

}

func trades(conn *grpc.ClientConn) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("-- coin check trades --")
	in := &bitco.TradesParams{Pair: "btc_jpy"}
	item, err := c.Trades(ctx, in)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(item)
	// b, err := json.Marshal(item)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(string(b))
	fmt.Println()
}

func orderBooks(conn *grpc.ClientConn) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("-- coin check order books --")
	in := &bitco.Empty{}
	item, err := c.OrderBooks(ctx, in)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(item)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	fmt.Println()
}

func RatePair(conn *grpc.ClientConn) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("-- coin check order books --")
	in := &bitco.RatePairParams{
		Pair: bitco.Btcjpy.String(),
	}
	item, err := c.RatePair(ctx, in)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(item)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	fmt.Println()
}

func TickHist(conn *grpc.ClientConn) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("-- coin check ticker history --")
	in := bitco.TickerHistParam{
		Limit: 10,
	}
	item, err := c.TickerHist(ctx, &in)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(item)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	fmt.Println()
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	ticker(conn)
	// trades(conn)
	orderBooks(conn)
	RatePair(conn)
	TickHist(conn)
	os.Exit(0)
}
