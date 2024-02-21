# wallet-service

## Installation

You need `sqlite3` installed and `curl` if you want to try the examples from bellow. Then you can either run with go run
or build the binary

```shell
go run main.go
```

or

```shell
go build
./wallet-service
```

A Dockerfile is also provided to run the project:

```shell
git clone git@github.com:zuzuleinen/wallet-service.git
cd wallet-service
docker build -t walletservice .
docker run -p 8080:8080 walletservice
```

## Using the service

Check Health of the Service:

```shell
curl http://localhost:8080/health
```

Create a wallet with amount `100` cents for user with id `andrei`:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 100}' http://localhost:8080/wallet/andrei
```

Query the current state of a wallet:

```shell
curl http://localhost:8080/wallet/andrei
```

Add funds to a wallet. Make sure the `reference` field is unique:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 10,"reference":"wonbet-1"}' http://localhost:8080/add-funds/andrei
```

Remove funds from wallet. Make sure the `reference` field is unique:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 50,"reference":"lostbet-2"}' http://localhost:8080/remove-funds/andrei
```