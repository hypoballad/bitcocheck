package bitcocheck

import (
	"reflect"
	"testing"
)

func TestRatePaircc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
		pair Pair
	}
	tests := []struct {
		name    string
		args    args
		want    RatePairItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "rate pair test",
			args:    args{conf: conf, pair: Btcjpy},
			want:    RatePairItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RatePaircc(tt.args.conf, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatePaircc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RatePaircc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarketBuycc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf   Config
		pair   Pair
		amount uint32
	}
	tests := []struct {
		name    string
		args    args
		want    MarketItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "market buy test",
			args: args{
				conf:   conf,
				pair:   Btcjpy,
				amount: 500,
			},
			want:    MarketItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarketBuycc(tt.args.conf, tt.args.pair, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarketBuycc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarketBuycc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExchangeOrdersOpenscc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		want    OrdersOpensItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "order opens test",
			args:    args{conf: conf},
			want:    OrdersOpensItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExchangeOrdersOpenscc(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExchangeOrdersOpenscc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExchangeOrdersOpenscc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExchangeOrdersTransactionscc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		want    OrdersTransactionsItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "orders transactions test",
			args:    args{conf: conf},
			want:    OrdersTransactionsItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExchangeOrdersTransactionscc(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExchangeOrdersTransactionscc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExchangeOrdersTransactionscc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountsBalancecc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		want    AccountsBalanceItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "accounts balance test",
			args:    args{conf: conf},
			want:    AccountsBalanceItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AccountsBalancecc(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountsBalancecc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountsBalancecc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccounts(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		want    AccountsItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "accounts test",
			args:    args{conf: conf},
			want:    AccountsItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Accountscc(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Accounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Accounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTickercc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		want    TickerItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "ticker test",
			args:    args{conf: conf},
			want:    TickerItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tickercc(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tickercc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tickercc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTradescc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
		pair Pair
	}
	tests := []struct {
		name    string
		args    args
		want    TradesItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "trade test",
			args:    args{conf: conf, pair: Btcjpy},
			want:    TradesItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tradescc(tt.args.conf, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tradescc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tradescc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderBookscc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		want    OrderBooksItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "order books test",
			args:    args{conf: conf},
			want:    OrderBooksItem{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OrderBookscc(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderBookscc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderBookscc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExchangeOrdersRatecc(t *testing.T) {
	conf, err := DecodeConfigToml("bitcocheck.toml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		conf        Config
		order       OrderType
		pair        Pair
		amountprice AmountPriceType
		value       string
	}
	tests := []struct {
		name    string
		args    args
		want    ExchangeOrdersRateItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "exchange order rate test",
			args: args{
				conf:        conf,
				order:       Buy,
				pair:        Btcjpy,
				amountprice: Price,
				value:       "10000",
			},
			want:    ExchangeOrdersRateItem{},
			wantErr: false,
		},
		{
			name: "exchange order rate test",
			args: args{
				conf:        conf,
				order:       Sell,
				amountprice: Price,
				pair:        Btcjpy,
				value:       "10000",
			},
			want:    ExchangeOrdersRateItem{},
			wantErr: false,
		},
		{
			name: "exchange order rate test",
			args: args{
				conf:        conf,
				order:       Buy,
				pair:        Btcjpy,
				amountprice: Amount,
				value:       "0.1",
			},
			want:    ExchangeOrdersRateItem{},
			wantErr: false,
		},
		{
			name: "exchange order rate test",
			args: args{
				conf:        conf,
				order:       Sell,
				amountprice: Amount,
				pair:        Btcjpy,
				value:       "0.1",
			},
			want:    ExchangeOrdersRateItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExchangeOrdersRatecc(tt.args.conf, tt.args.order, tt.args.pair, tt.args.amountprice, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExchangeOrdersRatecc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExchangeOrdersRatecc() = %v, want %v", got, tt.want)
			}
		})
	}
}
