package bitcocheck

import (
	"encoding/json"
	fmt "fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Main MainConfig `toml:"main"`
}

type MainConfig struct {
	Access string `toml:"access"`
	Secret string `toml:"secret"`
	Debug  bool   `toml:"debug"`
}

// DecodeConfigToml ...
func DecodeConfigToml(tomlfile string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(tomlfile, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

const CoincheckURL = "https://coincheck.com"

func targetAPI(api string) string {
	return fmt.Sprintf("%s%s", CoincheckURL, api)
}

// Tickercc You can get the latest information easily.
func Tickercc(conf Config) (TickerItem, error) {
	var tickerItem TickerItem
	url := targetAPI("/api/ticker")
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return tickerItem, err
	}
	if err := json.Unmarshal(jsonBlob, &tickerItem); err != nil {
		return tickerItem, err
	}
	return tickerItem, nil
}

type Pair int

const (
	Btcjpy Pair = iota
	Fctjpy
)

func (p Pair) String() string {
	return [...]string{"btc_jpy", "fct_jpy"}[p]
}

// Tradescc You can get the latest transaction history.
func Tradescc(conf Config, pair Pair) (TradesItem, error) {
	var item TradesItem
	url := targetAPI(fmt.Sprintf("/api/trades?pair=%s", pair.String()))
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return item, err
	}
	if err := json.Unmarshal(jsonBlob, &item); err != nil {
		return item, err
	}
	return item, nil
}

type OrderBooksItemIntermediate struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

// OrderBookscc Board information can be obtained.
func OrderBookscc(conf Config) (OrderBooksItem, error) {
	var item OrderBooksItem
	url := targetAPI("/api/order_books")
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return item, err
	}
	// fmt.Println(string(jsonBlob))
	var intermediate OrderBooksItemIntermediate
	if err := json.Unmarshal(jsonBlob, &intermediate); err != nil {
		return item, err
	}
	asksOrderArray := []*OrderArray{}
	for _, asks := range intermediate.Asks {
		ask := []string{}
		for _, a := range asks {
			ask = append(ask, a)
		}
		askArray := OrderArray{Items: ask}
		asksOrderArray = append(asksOrderArray, &askArray)
	}
	bidsOrderArray := []*OrderArray{}
	for _, bids := range intermediate.Bids {
		bid := []string{}
		for _, b := range bids {
			bid = append(bid, b)
		}
		bidArray := OrderArray{Items: bid}
		bidsOrderArray = append(bidsOrderArray, &bidArray)
	}
	item.Asks = asksOrderArray
	item.Bids = bidsOrderArray
	return item, nil
}

// OrderType Note method
type OrderType int

const (
	// Sell Limit order, spot trading, sell.
	Sell OrderType = iota
	// Buy Limit order, spot trading, buy.
	Buy
	// MarketBuy Market order Cash transaction Buy
	MarketBuy
	// MarketSell Market orders, spot trading, selling
	MarketSell
)

func (o OrderType) String() string {
	return [...]string{"sell", "buy", "market_buy", "market_sell"}[o]
}

type AmountPriceType int

const (
	Amount AmountPriceType = iota
	Price
)

func (a AmountPriceType) String() string {
	return [...]string{"amount", "price"}[a]
}

// ExchangeOrdersRatecc The rate is calculated based on the exchange's order.
func ExchangeOrdersRatecc(conf Config, order OrderType, pair Pair, amountprice AmountPriceType, value string) (ExchangeOrdersRateItem, error) {
	var exchangeOrdersRateItem ExchangeOrdersRateItem
	url := targetAPI(fmt.Sprintf("/api/exchange/orders/rate?order_type=%s&pair=%s&%s=%s", order.String(), pair.String(), amountprice, value))
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return exchangeOrdersRateItem, err
	}
	//fmt.Println(string(jsonBlob))
	if err := json.Unmarshal(jsonBlob, &exchangeOrdersRateItem); err != nil {
		return exchangeOrdersRateItem, err
	}
	// fmt.Println(exchangeOrdersRateItem)
	return exchangeOrdersRateItem, nil
}

// RatePaircc Get a dealership rate
func RatePaircc(conf Config, pair Pair) (RatePairItem, error) {
	var ratePairItem RatePairItem
	url := targetAPI(fmt.Sprintf("/api/rate/%s", pair.String()))
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return ratePairItem, err
	}
	if err := json.Unmarshal(jsonBlob, &ratePairItem); err != nil {
		return ratePairItem, err
	}
	return ratePairItem, nil
}

type MarketBuyPayload struct {
	Pair            string `json:"pair"`
	OrderType       string `json:"order_type"`
	MarketBuyAmount uint32 `json:"market_buy_amount"`
}

// MarketBuycc Market order Cash transaction Buy
func MarketBuycc(conf Config, pair Pair, amount uint32) (MarketItem, error) {
	var marketItem MarketItem
	url := targetAPI("/api/exchange/orders")
	payload := MarketBuyPayload{
		Pair:            pair.String(),
		OrderType:       MarketBuy.String(),
		MarketBuyAmount: amount,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return marketItem, err
	}
	body := string(payloadBytes)
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.PostRequest()
	if err != nil {
		return marketItem, err
	}
	fmt.Println(string(jsonBlob))
	if err := json.Unmarshal(jsonBlob, &marketItem); err != nil {
		return marketItem, err
	}
	fmt.Println(marketItem)
	return marketItem, nil
}

type MarketSellPayload struct {
	Pair      string `json:"pair"`
	OrderType string `json:"order_type"`
	Amount    uint32 `json:"amount"`
}

// MarketSellcc Market orders, spot trading, selling
func MarketSellcc(conf Config, pair Pair, amount uint32) (MarketItem, error) {
	var marketItem MarketItem
	url := targetAPI("/api/exchange/orders")
	payload := MarketSellPayload{
		Pair:      pair.String(),
		OrderType: MarketSell.String(),
		Amount:    amount,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return marketItem, err
	}
	body := string(payloadBytes)
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.PostRequest()
	if err != nil {
		return marketItem, err
	}
	fmt.Println(string(jsonBlob))
	if err := json.Unmarshal(jsonBlob, &marketItem); err != nil {
		return marketItem, err
	}
	fmt.Println(marketItem)
	return marketItem, nil
}

type LimitOrderPayload struct {
	Pair      string `json:"pair"`
	OrderType string `json:"order_type"`
	Rate      string `json:"rate"`
	Amount    string `json:"string"`
	// Positonid    int64  `json:"position_id,omitempy"`
	StopLossRate string `json"stop_loss_rate,omitempty"`
}

// LimitOrdercc Limit order, spot trading, buy.
func LimitOrdercc(conf Config, pair Pair, ordertype OrderType, rate, amount, stoplossrate string) (MarketItem, error) {
	var marketItem MarketItem
	url := targetAPI("/api/exchange/orders")
	// now := time.Now().Unix()
	payload := LimitOrderPayload{
		Pair:      pair.String(),
		OrderType: ordertype.String(),
		Rate:      rate,
		Amount:    amount,
		// Positonid:    now,
		StopLossRate: stoplossrate,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return marketItem, err
	}
	body := string(payloadBytes)
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.PostRequest()
	if err != nil {
		return marketItem, err
	}
	fmt.Println(string(jsonBlob))
	if err := json.Unmarshal(jsonBlob, &marketItem); err != nil {
		return marketItem, err
	}
	fmt.Println(marketItem)
	return marketItem, nil
}

// ExchangeOrdersOpenscc View a list of pending orders in your account.
func ExchangeOrdersOpenscc(conf Config) (OrdersOpensItem, error) {
	var item OrdersOpensItem
	url := targetAPI("/api/exchange/orders/opens")
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return item, err
	}
	if err := json.Unmarshal(jsonBlob, &item); err != nil {
		return item, err
	}
	return item, nil
}

func DeleteExchangeOrdercc(conf Config, id uint32) (DeleteOrderItem, error) {
	var item DeleteOrderItem
	url := targetAPI(fmt.Sprintf("/api/exchange/orders/%d", id))
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Delete()
	if err != nil {
		return item, err
	}
	fmt.Println(jsonBlob)
	if err := json.Unmarshal(jsonBlob, &item); err != nil {
		return item, err
	}
	fmt.Println(item)
	return item, nil
}

// ExchangeOrdersTransactionscc You can see your recent transaction history.
func ExchangeOrdersTransactionscc(conf Config) (OrdersTransactionsItem, error) {
	var item OrdersTransactionsItem
	url := targetAPI("/api/exchange/orders/transactions")
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return item, err
	}
	if err := json.Unmarshal(jsonBlob, &item); err != nil {
		return item, err
	}
	return item, nil
}

// AccountsBalancecc You can check the balance of your account.
func AccountsBalancecc(conf Config) (AccountsBalanceItem, error) {
	var item AccountsBalanceItem
	url := targetAPI("/api/accounts/balance")
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return item, err
	}
	if err := json.Unmarshal(jsonBlob, &item); err != nil {
		return item, err
	}
	return item, nil
}

// Accounts View your account information.
func Accountscc(conf Config) (AccountsItem, error) {
	var item AccountsItem
	url := targetAPI("/api/accounts")
	body := ""
	apiInfo := NewAPIInfo(conf.Main.Access, conf.Main.Secret, url, body, conf.Main.Debug)
	jsonBlob, err := apiInfo.Request()
	if err != nil {
		return item, err
	}
	// fmt.Println(string(jsonBlob))
	if err := json.Unmarshal(jsonBlob, &item); err != nil {
		return item, err
	}
	// fmt.Println(item)
	return item, nil
}
