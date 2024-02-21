# wallet-service

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

## Notes

- For simplicity, I assume the users already exists.
- Normally a user_id can be some UUID, but for ease of testing will just use regular strings.
- For managing funds I created 2 endpoints `/add-funds/{userId}` and `/remove-funds/{userId`. Depending on the business
  requirements and logic could also have something more concrete ie. `/place-bet/{userId}`
- I'm doing request validation in the HTTP handlers. For more complex requests could have also created some validators.
- // InitDatabase opens DB connection and also migrate for simplicity. In a real-case scenario I would migrate on a
  build step not when running the server
- We should have UNIQUE constraint on user_id and reference. for now I just concatenate them to solve the issue
- TODO. Maybe I can allow same transactions to be saved and when computing the balance I just ignore the ones already
  added?
- TODO: add proper logging
