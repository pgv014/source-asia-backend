# Source Asia Backend Assignment

## Overview

This project implements the required backend assignment for Source Asia using Go and in-memory storage.

The project contains:

1. Rate-limited API
2. Product catalog API with media support

The implementation focuses on:
- concurrency safety
- API design
- validation
- pagination
- performance-aware list responses.

---

# Tech Stack

- Go
- net/http
- In-memory data structures (maps + slices)
- sync.Mutex for concurrency safety

---

# Project Structure

```text
source-asia/
├── main.go
├── README.md
├── internal/
│   ├── api/
│   │   ├── request_handler.go
│   │   └── stats_handler.go
│   ├── products/
│   │   ├── handler.go
│   │   ├── model.go
│   │   └── store.go
│   └── ratelimit/
│       └── limiter.go
```

---

# How to Run

## Install dependencies

```bash
go mod tidy
```

## Run server

```bash
go run .
```

Server starts on:

```text
http://localhost:8080
```

---

# Part 1 — Rate Limited API

## POST /request

Accepts requests with per-user rate limiting.

### Request Body

```json
{
  "user_id": "alice",
  "payload": {
    "message": "hello"
  }
}
```

### Success Response

Status:
```http
201 Created
```

Example:
```json
{
  "status": "accepted",
  "user_id": "alice"
}
```

### Validation Errors

Returns:

```http
400 Bad Request
```

Cases:
- missing user_id
- empty user_id
- invalid JSON
- missing payload

### Rate Limiting

Implementation:
- rolling 1-minute window
- maximum 5 accepted requests per user

When exceeded:

```http
429 Too Many Requests
```

Example:
```json
{
  "error": "rate limit exceeded"
}
```

---

## GET /stats

Returns per-user request statistics.

### Example Response

```json
{
  "alice": {
    "accepted_requests": 5,
    "rejected_requests": 1
  }
}
```

Rejected requests are cumulative.

---

# Part 2 — Product Catalog API

## Product Model

```json
{
  "id": "prod_000001",
  "name": "Widget A",
  "sku": "SKU-001",
  "image_urls": [
    "https://cdn.example.com/products/sku-001/img-1.jpg"
  ],
  "video_urls": [
    "https://cdn.example.com/products/sku-001/demo.mp4"
  ]
}
```

---

## POST /products

Creates a product.

### Validation

- name required
- sku required
- sku must be unique
- only http/https URLs allowed
- max URL length: 2048
- max 20 URLs per request array

### Success Response

```http
201 Created
```

---

## GET /products

Returns paginated lightweight product list.

### Query Parameters

- limit
- offset

### Example

```text
GET /products?limit=5&offset=0
```

### Important Performance Design

This endpoint intentionally does NOT return:
- image_urls
- video_urls

Instead it returns:
- image_count
- video_count

This avoids loading and serializing large media arrays for list pages.

### Example Response

```json
[
  {
    "id": "prod_000001",
    "name": "Widget A",
    "sku": "SKU-001",
    "image_count": 2,
    "video_count": 1
  }
]
```

---

## GET /products/{id}

Returns full product details including all media URLs.

### Unknown Product

```http
404 Not Found
```

---

## POST /products/{id}/media

Appends media URLs to an existing product.

### Request Example

```json
{
  "image_urls": [
    "https://cdn.example.com/products/sku-001/new.jpg"
  ]
}
```

### Validation

At least one of:
- image_urls
- video_urls

must be present.

---

# Concurrency Safety

The project uses mutex locks (`sync.Mutex`) for:
- rate limiter state
- product storage

This ensures correctness under concurrent requests.

---

# Storage Design

Products are stored in-memory using:
- product map
- SKU index map

Media URLs are stored separately inside product objects.

No database is used.

---

# Production Limitations

Current implementation limitations:

- in-memory only
- data lost on restart
- single-instance deployment
- no distributed rate limiting
- no persistent storage
- no authentication

---

# Production Improvements

In production, this system could be improved using:
- PostgreSQL for persistence
- Redis for distributed rate limiting
- CDN-backed media storage
- cursor-based pagination
- background workers
- structured logging

---

# Example curl Commands

## Rate Limit Test

```bash
for i in $(seq 1 6); do
curl -s -X POST http://localhost:8080/request \
-H 'Content-Type: application/json' \
-d '{"user_id":"alice","payload":{"n":'$i'}}'
echo ""
done
```

---

## Stats

```bash
curl -s http://localhost:8080/stats
```

---

## Create Product

```bash
curl -s -X POST http://localhost:8080/products \
-H 'Content-Type: application/json' \
-d '{
"name":"Widget A",
"sku":"SKU-001",
"image_urls":["https://cdn.example.com/img1.jpg"],
"video_urls":["https://cdn.example.com/demo.mp4"]
}'
```

---

## List Products

```bash
curl -s "http://localhost:8080/products?limit=5&offset=0"
```

---

## Product Detail

```bash
curl -s http://localhost:8080/products/prod_000001
```

---



All code was reviewed, modified, tested, and validated manually.
