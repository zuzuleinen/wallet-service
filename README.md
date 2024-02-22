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

To run the tests:

```shell
go test ./...
```

A Dockerfile is also provided to run the project:

```shell
git clone git@github.com:zuzuleinen/wallet-service.git
cd wallet-service
docker build -t walletservice .
docker run -p 8080:8080 walletservice
```

## Using the service

For simplicity, user ids are simple strings, like `andrei`. Amounts are in cents and always a positive value.

**Check Health of the Service**

```shell
curl http://localhost:8080/health
```

**Create a Wallet**

Create a wallet with amount `100` cents for user with id `andrei`:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 100}' http://localhost:8080/wallet/andrei
```

**Query the current state of a wallet**

```shell
curl http://localhost:8080/wallet/andrei
```

**Add funds to a wallet**

Make sure the `reference` field is unique:
A `reference` contains information about the context in which this amount is added/removed. For
example: `wonbet-111`, `witraw-222`, where the first part is the event and second it a unique id. I don't do any
validation for the format, but good when debugging.

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 10,"reference":"wonbet-1"}' http://localhost:8080/add-funds/andrei
```

**Remove funds from wallet**

Make sure the `reference` field is unique:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 50,"reference":"lostbet-2"}' http://localhost:8080/remove-funds/andrei
```