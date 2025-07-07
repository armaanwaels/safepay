# Safepay Backend API

A simple Go-based HTTP API for managing invoices and payments, with SQLite migrations and clean code structure.

##  Features

- RESTful API endpoints for creating and retrieving invoices, and processing/viewing payments
- Data stored in SQLite using automatic schema migration
- Easy to run locally, no external dependencies except Go and SQLite

---

## Project Structure

```
.
├── main.go
├── handlers/
│   ├── invoice.go
│   └── payment.go
├── models/
│   ├── invoice.go
│   └── payment.go
├── db/
│   └── schema.sql
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```



**Requirements:**

- Go 1.18+ installed
- SQLite3 (CLI or library)

## How to Run

1. **Install dependencies** (from the project root):
```
    go mod tidy
```
- Run the server:
```
    go run main.go
```
The server will start at:
http://localhost:8080

On first run, the database (safepay.db) and tables will be created automatically.

Reset the Database:
```
rm safepay.db
go run main.go
```
## Example Usage (with curl)

- Create an Invoice:
```
curl -X POST http://localhost:8080/invoices \
  -H "Content-Type: application/json" \
  -d '{"merchant_id":"abc123","customer_id":"cus789","amount":2500,"currency":"USD","description":"Test invoice"}'
```

- Fetch an Invoice:
```
curl http://localhost:8080/invoices/<invoice_id>
```

- rocess a Payment:
```
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -d '{"invoice_id":"<invoice_id>", "method":"card", "details":{"card_number":"4111 1111 1111 1111", "expiry":"12/25", "cvv":"123"}}'
```
- Fetch Payment Status:
```
curl http://localhost:8080/payments/<invoice_id>
Replace <invoice_id> with the real invoice ID returned from the create invoice step.
```

## API Endpoints

- Invoices:

  POST /invoices
  Create a new invoice:
  ```
  {
    "merchant_id": "abc123",
    "customer_id": "cus789",
    "amount": 2500,
    "currency": "USD",
    "description": "Test invoice"
  }
  ```
  GET /invoices/{id}
  Get a specific invoice by ID

- Payments:

    POST /payments
    Process a payment for an invoice:
    
      {
        "invoice_id": "inv_xxxxxxxx",
        "method": "card",
        "details": {
          "card_number": "4111 1111 1111 1111",
          "expiry": "12/25",
          "cvv": "123"
        }
      }
      
    GET /payments/{invoice_id}
    Get the status of payment for an invoice

## Running Tests

Unit tests are included for core handlers.

go test ./handlers

## Armaan Waels - Golang API assignment – Safepay (2025)

