# Bank system in Go

## Simple API

https://go-bank-service.herokuapp.com/

### Using

```http request
GET /cards
Content-Type: application/json

{
  "clientId": 1
}
```

```http request
GET /transactions
Content-Type: application/json

{
  "cardId": 1
}
```

```http request
GET /most-expensive
Content-Type: application/json

{
  "cardId": 1
}
```

```http request
GET /most-popular
Content-Type: application/json

{
  "cardId": 1
}
```
