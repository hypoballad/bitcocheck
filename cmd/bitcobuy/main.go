package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	bitco "github.com/hypoballad/bitcocheck"
	"github.com/rs/xid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:50051", "server address")
var actualMode = flag.Bool("actual", false, "actual mode")
var debugMode = flag.Bool("debug", false, "mode debug")
var commandName = flag.String("c", "", "")
var dbFile = flag.String("db", "bitcobuy.db", "")

var conf *sqlite3.Conn

const OrderInfo = `create table if not exists order_info (
	id text PRIMARY_KEY,
	order_id int NOT NULL,
	order_type text NOT NULL,
	ts timestamp NOT NULL,
	btc text NOT NULL,
	yen text NOT NULL,
	item json NOT NULL
)
`

const TradeHist = `create table if not exists trade_hist (
	id text PRIMARY_KEY,
	ts timestamp NOT NULL,
	btc text not NULL,
	yen text NOT NULL
)
`

func createSQL(conn *sqlite3.Conn) error {
	for _, stmt := range []string{OrderInfo, TradeHist} {
		if err := conn.Exec(stmt); err != nil {
			return errors.New(fmt.Sprintf("%v, %s", err, stmt))
		}
	}
	return nil
}

type Hist struct {
	ID  string
	Ts  string
	Btc string
	Yen string
}

func FindTradeHist(conn *sqlite3.Conn) ([]Hist, error) {
	hists := []Hist{}
	stmt, err := conn.Prepare(`select id, ts, btc, yen from trade_hist order by ts asc`)
	if err != nil {
		return hists, err
	}
	defer stmt.Close()
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return hists, err
		}
		if !hasRow {
			break
		}
		var id string
		var ts string
		var btc string
		var yen string
		if err := stmt.Scan(&id, &ts, &btc, &yen); err != nil {
			return hists, err
		}
		hists = append(hists, Hist{ID: id, Ts: ts, Btc: btc, Yen: yen})
	}
	return hists, nil
}

func SaveTradeHist(conn *sqlite3.Conn, btc, yen string) error {
	guid := xid.New()
	if err := conn.Begin(); err != nil {
		return err
	}
	stmt, err := conn.Prepare(`insert into trade_hist values (?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	tm := time.Now()

	if err := stmt.Exec(guid.String(), tm.Format("2006-01-02 15:04:05"), btc, yen); err != nil {
		return err
	}
	if err := conn.Commit(); err != nil {
		return err
	}
	return nil
}

type Order struct {
	ID        string
	OredrID   uint32
	OrderType string
	Ts        string
	Btc       string
	Yen       string
	Item      string
}

func FindBuyList(conn *sqlite3.Conn) ([]Order, error) {
	orders := []Order{}
	stmt, err := conn.Prepare(`select id, order_id, order_type, ts, btc, yen, item from order_info order by ts asc`)
	if err != nil {
		return orders, err
	}
	defer stmt.Close()
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return orders, err
		}
		if !hasRow {
			break
		}
		var id string
		var orderid int
		var ordertype string
		var ts string
		var btc string
		var yen string
		var item string
		if err := stmt.Scan(&id, &orderid, &ordertype, &ts, &btc, &yen, &item); err != nil {
			return orders, err
		}
		orders = append(orders, Order{ID: id, OredrID: uint32(orderid), OrderType: ordertype, Ts: ts, Btc: btc, Yen: yen, Item: item})
	}
	return orders, nil
}

func DelBuyInfo(conn *sqlite3.Conn, id string) error {
	if err := conn.Begin(); err != nil {
		return err
	}
	stmt, err := conn.Prepare(`delete order_info where id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Exec(id); err != nil {
		return err
	}
	if err := conn.Commit(); err != nil {
		return err
	}
	return nil
}

func SaveBuyInfo(conn *sqlite3.Conn, orderid int, ordertype, btc, yen string, item *bitco.MarketItem) error {
	guid := xid.New()
	if err := conn.Begin(); err != nil {
		return err
	}
	stmt, err := conn.Prepare(`insert into order_info values (?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	tm := time.Now()
	bstr, err := json.Marshal(item)
	if err != nil {
		return err
	}
	if err := stmt.Exec(guid.String(), orderid, ordertype, tm.Format("2006-01-02 15:04:05"), btc, yen, string(bstr)); err != nil {
		return err
	}
	if err := conn.Commit(); err != nil {
		return err
	}
	return nil
}

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

func TotalAssets(addr string, debug bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// account balance
	balance, err := AccountsBalance(conn)
	if err != nil {
		fmt.Println("accounts balance error: ", err)
		return
	}
	yen, err := strconv.ParseFloat(balance.Jpy, 32)
	if err != nil {
		fmt.Println("jpy convert error: ", err)
		return
	}
	// debugJson(balance)
	btc, err := strconv.ParseFloat(balance.Btc, 32)
	if err != nil {
		fmt.Println("btc convert error:", err)
		return
	}
	salesrate, err := SalesRate(conn)
	if err != nil {
		log.Println("sales rate error:", err)
		return
	}
	rate, err := strconv.ParseFloat(salesrate.Rate, 32)
	if err != nil {
		log.Println("sales rate convert error:", err)
		return
	}
	fmt.Println("== 総資産 ==")
	fmt.Printf("資金: %s 円\n", humanizeYen(balance.Jpy))
	fmt.Printf("BTC:  %f BTC (%s円)\n", btc, humanizeYen(fmt.Sprintf("%f", rate*btc)))
	fmt.Printf("総額: %s円\n", humanizeYen(fmt.Sprintf("%f", rate*btc+yen)))
	fmt.Println()
}

func humanizeYen(yen string) string {
	humanize := ""
	f, err := strconv.ParseFloat(yen, 32)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	p := message.NewPrinter(language.English)
	humanize = p.Sprintf("%d", int(f))
	return humanize
}

func SuggestBuy(addr string, debug bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	balance, err := AccountsBalance(conn)
	if err != nil {
		fmt.Println("accounts balance error: ", err)
		return
	}
	buyrate, err := BuyRateBtc(conn, balance.Jpy)
	if err != nil {
		log.Println("buy rate error:", err)
	}
	// debugJson(item)
	salesrate, err := SalesRate(conn)
	if err != nil {
		log.Println("sales rate error:", err)
	}
	fmt.Println("== 買いレート ==")
	fmt.Printf("レート: %s 円(1btc)\n", humanizeYen(salesrate.Rate))
	fmt.Printf("買値: %s 円(1btc)\n", humanizeYen(buyrate.Rate))
	fmt.Printf("%s円 : %sbtc\n", humanizeYen(buyrate.Price), buyrate.Amount)
	fmt.Println()
}

func SuggestSell(addr string, debug bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	balance, err := AccountsBalance(conn)
	if err != nil {
		fmt.Println("accounts balance error: ", err)
		return
	}
	sellrate, err := SellRateBtc(conn, balance.Btc)
	if err != nil {
		log.Println("sell rate error:", err)
		return
	}
	// debugJson(item)
	salesrate, err := SalesRate(conn)
	if err != nil {
		log.Println("sales rate error:", err)
		return
	}
	fmt.Println("== 売りレート ==")
	fmt.Printf("レート: %s 円(1btc)\n", humanizeYen(salesrate.Rate))
	fmt.Printf("売り値: %s 円(1btc)\n", humanizeYen(sellrate.Rate))
	fmt.Printf("%s円 : %sbtc\n", humanizeYen(sellrate.Price), sellrate.Amount)
	fmt.Println()
}

func ExchangeOrdersOpens(conn *grpc.ClientConn) (*bitco.OrdersOpensItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := &bitco.Empty{}
	item, err := c.ExchangeOrdersOpens(ctx, in)
	if err != nil {
		return item, err
	}
	return item, nil
}

func Pendings(addr string, debug bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return
	}
	defer conn.Close()
	items, err := ExchangeOrdersOpens(conn)
	if err != nil {
		log.Println("exchange order open error:", err)
		return
	}
	fmt.Println("== 未決済一覧 ==")
	for _, item := range items.Orders {
		fmt.Printf("ID: %d\n", item.Id)
		fmt.Printf("売買: %d\n", item.OrderType)
		fmt.Printf("レート: %d\n", item.Rate)
		fmt.Printf("量: %s\n", item.PendingAmount)
		fmt.Println()
	}
}

func DeleteExchangeOrder(conn *grpc.ClientConn, id uint32) (uint32, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := &bitco.DeleteOrderParam{Id: id}
	item, err := c.DeleteExchangeOrder(ctx, in)
	if err != nil {
		return 0, err
	}
	return item.Id, nil
}

func CancelOrder(sqlcon *sqlite3.Conn, addr string, debug bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return
	}
	defer conn.Close()
	orders, err := FindBuyList(sqlcon)
	if err != nil {
		log.Println("find buy list error:", err)
		return
	}
	fmt.Println("== 注文キャンセル ==")
	if len(orders) == 0 {
		fmt.Println("注文はありません")
		return
	}
	id, err := DeleteExchangeOrder(conn, orders[0].OredrID)
	if err != nil {
		fmt.Println("注文キャンセルエラー:", err)
		return
	}
	fmt.Printf("注文をキャンセルしました: %d\n", id)
	debugJson(orders[0])

}

func LimitBuy(conn *grpc.ClientConn, in *bitco.LimitOrderParams) (*bitco.MarketItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	item, err := c.LimitBuy(ctx, in)
	if err != nil {
		return item, err
	}
	return item, nil
}

func BuyOrder(sqlcon *sqlite3.Conn, addr string, actual bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	orders, err := FindBuyList(sqlcon)
	if err != nil {
		log.Println("find buy list error:", err)
		return
	}
	fmt.Println("== 買い注文 ==")
	if len(orders) != 0 {
		fmt.Println("すでにポジションを持ってます")
		item := orders[0]
		// fmt.Println(item)
		//fmt.Printf("買値: %s 円(1btc)\n", humanizeYen(item.Rate))
		fmt.Printf("注文番号:%d\n", item.OredrID)
		fmt.Printf("%s円 : %sbtc\n", humanizeYen(item.Yen), item.Btc)
		fmt.Println()
		return
	}

	balance, err := AccountsBalance(conn)
	if err != nil {
		fmt.Println("accounts balance error: ", err)
		return
	}
	yen, err := strconv.ParseFloat(balance.Jpy, 32)
	if err != nil {
		fmt.Println("jpy convert error: ", err)
		return
	}
	if int(yen) == 0 {
		fmt.Println("軍資金がゼロです")
		return
	}
	buyrate, err := BuyRateBtc(conn, balance.Jpy)
	if err != nil {
		log.Println("buy rate error:", err)
		return
	}

	// debugJson(item)
	salesrate, err := SalesRate(conn)
	if err != nil {
		log.Println("sales rate error:", err)
	}

	in := bitco.LimitOrderParams{}
	in.Pair = bitco.Btcjpy.String()
	in.Rate = buyrate.Rate
	in.Amount = buyrate.Amount
	// debugJson(in)
	var item *bitco.MarketItem
	if actual {
		// LimitBuy
	} else {
		now := time.Now()
		item = &bitco.MarketItem{
			Success:      "true",
			Id:           uint64(now.Unix()),
			Rate:         buyrate.Rate,
			Amount:       buyrate.Amount,
			OrderType:    bitco.Buy.String(),
			StopLossRate: "",
			Pair:         bitco.Btcjpy.String(),
			CreatedAt:    now.Format("2006-01-02 15:04:05"),
		}
	}

	if err := SaveBuyInfo(sqlcon, int(item.Id), item.OrderType, item.Amount, buyrate.Price, item); err != nil {
		log.Println("save buy info error:", err)
		return
	}

	fmt.Printf("レート: %s 円(1btc)\n", humanizeYen(salesrate.Rate))
	fmt.Printf("買値: %s 円(1btc)\n", humanizeYen(item.Rate))
	fmt.Printf("%s円 : %sbtc\n", humanizeYen(buyrate.Price), item.Amount)
	fmt.Println()
}

func LimitSell(conn *grpc.ClientConn, in *bitco.LimitOrderParams) (*bitco.MarketItem, error) {
	c := bitco.NewCoincheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	item, err := c.LimitSell(ctx, in)
	if err != nil {
		return item, err
	}
	return item, nil
}

func SellOrder(sqlcon *sqlite3.Conn, addr string, actual bool) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	orders, err := FindBuyList(sqlcon)
	if err != nil {
		log.Println("find buy list error:", err)
		return
	}
	fmt.Println("== 売り注文 ==")
	if len(orders) > 0 {
		fmt.Println("ポジションはありません")
		return
	}

	salesrate, err := SalesRate(conn)
	if err != nil {
		log.Println("sales rate error:", err)
	}

	sellrate, err := SellRateBtc(conn, orders[0].Btc)
	if err != nil {
		log.Println("sell rate error:", err)
		return
	}

	in := bitco.LimitOrderParams{}
	in.Pair = bitco.Btcjpy.String()
	in.Rate = sellrate.Rate
	in.Amount = sellrate.Amount

	var item *bitco.MarketItem
	if actual {
		// LimitSell
	} else {
		now := time.Now()
		item = &bitco.MarketItem{
			Success:      "true",
			Id:           uint64(now.Unix()),
			Rate:         sellrate.Rate,
			Amount:       sellrate.Amount,
			OrderType:    bitco.Sell.String(),
			StopLossRate: "",
			Pair:         bitco.Btcjpy.String(),
			CreatedAt:    now.Format("2006-01-02 15:04:05"),
		}
	}

}

func main() {
	flag.Parse()
	conn, err := sqlite3.Open(*dbFile)
	if err != nil {
		log.Printf("sqlite3 connection error: %v\n", err)
		os.Exit(0)
	}
	if err := createSQL(conn); err != nil {
		log.Printf("create sql error: %v\n", err)
		os.Exit(0)
	}
	switch *commandName {
	case "assets":
		TotalAssets(*addr, *debugMode)
	case "buysuggest":
		SuggestBuy(*addr, *debugMode)
	case "sellsuggest":
		SuggestSell(*addr, *debugMode)
	case "pending":
		Pendings(*addr, *debugMode)
	case "cancel":
		CancelOrder(conn, *addr, *debugMode)
	case "limitbuy":
		BuyOrder(conn, *addr, *debugMode)
	default:
		log.Printf("command not found: %s\n", *commandName)
		os.Exit(0)
	}
	os.Exit(0)
}
