# bitcocheck

A sample GRPC server and client that makes a COINCHECK API request

## How to build GRPC

```
protoc -I bitcocheck/ bitcocheck/bitcocheck.proto --go_out=plugins=grpc:bitcocheck
```

## Example of a configuration file


Access keys and secret keys can be created [here](https://coincheck.com/ja/api_settings).

```
[main]
access = "CoinCheck API access key"
secret = "CoinCheck API secret key"
debug = false
```

## How to build bitcocheck command

```
go build -o butcocheck cmd/bitcocheck/main.go
```

## How to perform the bitcocheck

```
./bitcocheck -conf config.toml
```