
# Safepay Backend API

A simple Go-based HTTP API for managing invoices and payments, with SQLite migrations and clean code structure.  
Now fully Dockerized for easy deployment 🚀

---

## Features

- RESTful API endpoints for creating and retrieving invoices, and processing/viewing payments
- Data stored in SQLite using automatic schema migration
- Run locally with Go **or** in Docker using `docker compose`
- Persistent database storage via Docker volumes

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
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── README.md
└── .gitignore

````

---

## Requirements

### Option 1: Run Locally
- Go **1.24+** installed  
- SQLite3 (CLI or library)

### Option 2: Run with Docker (recommended)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed

---

## Running the API

### 🔹 Local (without Docker)
```bash
# Install dependencies
go mod tidy

# Run the server
go run main.go
````

The server will start at:
👉 [http://localhost:8080](http://localhost:8080)

Reset the database:

```bash
rm safepay.db
go run main.go
```

---

### 🔹 Docker (with Compose)

1. Build & run in detached mode:

```bash
docker compose up -d --build
```

2. View logs:

```bash
docker compose logs -f
```

3. Stop containers:

```bash
docker compose down
```

SQLite database persists in a Docker volume (`safepay_data`).

---

## Example Usage (with curl)

* **Create an Invoice**

```bash
curl -X POST http://localhost:8080/invoices \
  -H "Content-Type: application/json" \
  -d '{"merchant_id":"abc123","customer_id":"cus789","amount":2500,"currency":"USD","description":"Test invoice"}'
```

* **Fetch an Invoice**

```bash
curl http://localhost:8080/invoices/<invoice_id>
```

* **Process a Payment**

```bash
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -d '{"invoice_id":"<invoice_id>", "method":"card", "details":{"card_number":"4111 1111 1111 1111", "expiry":"12/25", "cvv":"123"}}'
```

* **Fetch Payment Status**

```bash
curl http://localhost:8080/payments/<invoice_id>
```

Replace `<invoice_id>` with the real invoice ID returned from the create step.

---

## API Endpoints

### Invoices

* **POST /invoices** – Create a new invoice
* **GET /invoices/{id}** – Fetch an invoice by ID

### Payments

* **POST /payments** – Process a payment for an invoice
* **GET /payments/{invoice\_id}** – Fetch payment status by invoice

---

## Running Tests

Unit tests are included for core handlers:

```bash
go test ./handlers
```

---

## Author

**Armaan Waels**
📌 Golang API Assignment – Safepay (2025)

---
