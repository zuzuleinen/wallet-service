# wallet-service

### Running with a local Pulsar cluster

1. Start Pulsar cluster locally:
```shell
sudo mkdir -p ./data/zookeeper ./data/bookkeeper
# this step might not be necessary on other than Linux platforms
sudo chown 10000 -R data
docker-compose up -d
```

2. Start the server locally:
```shell
go run cmd/rest-api/main.go
```

3. Start the background consumer responsible for persisting messages from Pulsar into db:
```shell
go run cmd/consumer/main.go
```
Or, you can go with the much faster batch consumer:
```shell
go run cmd/batch-consumer/main.go
```

Benchmark using [hey](https://github.com/rakyll/hey):

```shell
hey -n 100000 -m POST  -H "Content-Type: application/json" -d '{"amount": 10,"reference":"wonbet-1"}' http://localhost:8081/add-funds/andrei
```

On my machine last test was at `6726.2377` requests/sec.

## Run the tests

```shell
go test ./...
```

## Using the service

For simplicity, user ids are simple strings, like `andrei`. Amounts are in cents and always a positive value.

**Check Health of the Service**

```shell
curl http://localhost:8081/health
```

**Create a Wallet**

Create a wallet with amount `100` cents for user with id `andrei`:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 100}' http://localhost:8081/wallet/andrei
```

**Query the current state of a wallet**

```shell
curl http://localhost:8081/wallet/andrei
```

**Add funds to a wallet**

Make sure the `reference` field is unique:
A `reference` contains information about the context in which this amount is added/removed. For
example: `wonbet-111`, `witraw-222`, where the first part is the event and second it a unique id. I don't do any
validation for the format, but good when debugging.

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 10,"reference":"wonbet-1"}' http://localhost:8081/add-funds/andrei
```

**Remove funds from wallet**

Make sure the `reference` field is unique:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"amount": 50,"reference":"lostbet-2"}' http://localhost:8081/remove-funds/andrei
```