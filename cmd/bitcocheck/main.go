package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	bitco "github.com/hypoballad/bitcocheck"
	"github.com/robfig/cron/v3"
	"github.com/rs/xid"
	"google.golang.org/grpc"
)

const TickHist = `create table if not exists tickhist (
	id text PRIMARY_kEY,
	ts timestamp NOT NULL,
	last real NOT NULL,
	bid real NOT NULL,
	ask real NOT NULL,
	high real NOT NULL,
	low real NOT NULL,
	volume real NOT NULL)`

var addr = flag.String("addr", ":50051", "server address")
var configpath = flag.String("conf", "bitcocheck.toml", "config file name")
var dbFile = flag.String("db", "bitcocheck.db", "sqlite3 db file name")

var conf bitco.Config
var conn *sqlite3.Conn

type server struct {
	bitco.UnimplementedCoincheckServer
}

func (s server) Ticker(ctx context.Context, in *bitco.Empty) (*bitco.TickerItem, error) {
	var item bitco.TickerItem
	item, err := bitco.Tickercc(conf)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) Trades(ctx context.Context, in *bitco.TradesParams) (*bitco.TradesItem, error) {
	var item bitco.TradesItem
	var pair bitco.Pair
	switch in.Pair {
	case "ftc_jpy":
		pair = bitco.Fctjpy
	default:
		pair = bitco.Btcjpy
	}
	item, err := bitco.Tradescc(conf, pair)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) OrderBooks(ctx context.Context, in *bitco.Empty) (*bitco.OrderBooksItem, error) {
	var item bitco.OrderBooksItem
	item, err := bitco.OrderBookscc(conf)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) ExchangeOrdersRate(ctx context.Context, in *bitco.ExchangeOrdersRateParam) (*bitco.ExchangeOrdersRateItem, error) {
	var item bitco.ExchangeOrdersRateItem
	var orderType bitco.OrderType
	switch in.OrderType {
	case "sell":
		orderType = bitco.Sell
	default:
		orderType = bitco.Buy
	}
	var pair bitco.Pair
	switch in.Pair {
	case "ftc_jpy":
		pair = bitco.Fctjpy
	default:
		pair = bitco.Btcjpy
	}
	var amountPrice bitco.AmountPriceType
	switch in.Amountprice {
	case "amount":
		amountPrice = bitco.Amount
	default:
		amountPrice = bitco.Price
	}
	item, err := bitco.ExchangeOrdersRatecc(conf, orderType, pair, amountPrice, in.Value)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) RatePair(ctx context.Context, in *bitco.RatePairParams) (*bitco.RatePairItem, error) {
	var item bitco.RatePairItem
	var pair bitco.Pair
	switch in.Pair {
	case "ftc_jpy":
		pair = bitco.Fctjpy
	default:
		pair = bitco.Btcjpy
	}
	item, err := bitco.RatePaircc(conf, pair)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) LimitBuy(ctx context.Context, in *bitco.LimitOrderParams) (*bitco.MarketItem, error) {
	var item bitco.MarketItem
	var pair bitco.Pair
	switch in.Pair {
	case "ftc_jpy":
		pair = bitco.Fctjpy
	default:
		pair = bitco.Btcjpy
	}
	item, err := bitco.LimitOrdercc(conf, pair, bitco.Buy, in.Rate, in.Amount, in.StopLossRate)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) LimitSell(ctx context.Context, in *bitco.LimitOrderParams) (*bitco.MarketItem, error) {
	var item bitco.MarketItem
	var pair bitco.Pair
	switch in.Pair {
	case "ftc_jpy":
		pair = bitco.Fctjpy
	default:
		pair = bitco.Btcjpy
	}
	item, err := bitco.LimitOrdercc(conf, pair, bitco.Sell, in.Rate, in.Amount, in.StopLossRate)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) MarketBuy(ctx context.Context, in *bitco.MarketBuyParams) (*bitco.MarketItem, error) {
	var item bitco.MarketItem
	var pair bitco.Pair
	switch in.Pair {
	case "ftc_jpy":
		pair = bitco.Fctjpy
	default:
		pair = bitco.Btcjpy
	}
	item, err := bitco.MarketBuycc(conf, pair, in.MarketBuyAmount)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) MarketSell(ctx context.Context, in *bitco.MarketSellParam) (*bitco.MarketItem, error) {
	var item bitco.MarketItem
	var pair bitco.Pair
	switch in.Pair {
	case "ftc_jpy":
		pair = bitco.Fctjpy
	default:
		pair = bitco.Btcjpy
	}
	item, err := bitco.MarketSellcc(conf, pair, in.Amount)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) ExchangeOrdersOpens(ctx context.Context, in *bitco.Empty) (*bitco.OrdersOpensItem, error) {
	var item bitco.OrdersOpensItem
	item, err := bitco.ExchangeOrdersOpenscc(conf)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) DeleteExchangeOrer(ctx context.Context, in *bitco.DeleteOrderParam) (*bitco.DeleteOrderItem, error) {
	var item bitco.DeleteOrderItem
	item, err := bitco.DeleteExchangeOrdercc(conf, in.Id)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) ExchangeOrdersTransactions(ctx context.Context, in *bitco.Empty) (*bitco.OrdersTransactionsItem, error) {
	var item bitco.OrdersTransactionsItem
	item, err := bitco.ExchangeOrdersTransactionscc(conf)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) AccountsBalance(ctx context.Context, in *bitco.Empty) (*bitco.AccountsBalanceItem, error) {
	var item bitco.AccountsBalanceItem
	item, err := bitco.AccountsBalancecc(conf)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (s server) Accounts(ctx context.Context, in *bitco.Empty) (*bitco.AccountsItem, error) {
	var item bitco.AccountsItem
	//log.Println("accounts")
	item, err := bitco.Accountscc(conf)
	if err != nil {
		return &item, err
	}
	//fmt.Println(item)
	return &item, nil
}

func (s server) TickerHist(ctx context.Context, in *bitco.TickerHistParam) (*bitco.TickerHistItem, error) {
	var item bitco.TickerHistItem
	stmt, err := conn.Prepare(`select ts, last, bid, ask, high, low, volume from tickhist order by ts desc limit ?`)
	if err != nil {
		return &item, err
	}
	defer stmt.Close()
	limit := 1000
	if in.Limit > 0 {
		limit = int(in.Limit)
	}
	if err := stmt.Bind(limit); err != nil {
		return &item, err
	}
	result := []*bitco.TickerItem{}
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return &item, err
		}
		if !hasRow {
			break
		}
		var ts string
		var last float64
		var bid float64
		var ask float64
		var high float64
		var low float64
		var volume float64
		if err := stmt.Scan(&ts, &last, &bid, &ask, &high, &low, &volume); err != nil {
			return &item, err
		}
		tm, err := time.Parse("2006-01-02 15:04:05", ts)
		if err != nil {
			return &item, err
		}
		ticker := bitco.TickerItem{
			Timestamp: uint64(tm.Unix()),
			Last:      float32(last),
			Bid:       float32(bid),
			Ask:       float32(ask),
			High:      float32(high),
			Low:       float32(low),
			Volume:    float32(volume),
		}
		result = append(result, &ticker)
	}
	item.Tickeritem = result
	return &item, nil
}

func createSQL(conn *sqlite3.Conn) error {
	for _, stmt := range []string{TickHist} {
		if err := conn.Exec(stmt); err != nil {
			return errors.New(fmt.Sprintf("%v, %s", err, stmt))
		}
	}
	return nil
}

func job(conn *sqlite3.Conn, conf bitco.Config) error {
	item, err := bitco.Tickercc(conf)
	if err != nil {
		return err
	}
	guid := xid.New()
	if err := conn.Begin(); err != nil {
		return err
	}
	stmt, err := conn.Prepare(`insert into tickhist values (?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	tm := time.Unix(int64(item.Timestamp), 0)
	if err := stmt.Exec(guid.String(), tm.Format("2006-01-02 15:04:05"), float64(item.Last), float64(item.Bid), float64(item.Ask), float64(item.High), float64(item.Low), float64(item.Volume)); err != nil {
		return err
	}
	if err := conn.Commit(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	var err error
	conf, err = bitco.DecodeConfigToml(*configpath)
	if err != nil {
		log.Fatalln("config read error:", err)
	}
	conn, err = sqlite3.Open(*dbFile)
	if err != nil {
		log.Fatalln("sqlite3 connection error:", err)
	}
	if err := createSQL(conn); err != nil {
		log.Fatalln("create sql error: ", err)
	}
	if err := job(conn, conf); err != nil {
		log.Fatalln("job error:", err)
	}
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		if err := job(conn, conf); err != nil {
			log.Printf("job error %v\n", err)
		}
	})
	c.Start()
	var lis net.Listener
	lis, err = net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalln("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("listen to %s\n", *addr)
	bitco.RegisterCoincheckServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
