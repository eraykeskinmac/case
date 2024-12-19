# Invoice Management API - Case Study

A **RESTful API** service for managing invoices built with **Go**, **Fiber**, **GORM**, and **PostgreSQL**.

## 🚀 Tech Stack

- **Go** 1.23
- **Fiber** (Web Framework)
- **GORM** (ORM)
- **PostgreSQL**
- **Swagger** (API Documentation)
- **Docker & Docker Compose**

---

## 📑 API Endpoints

### Health Check

**`GET /api/health`**
Returns the health status of the API and its components.

### Invoices

#### List Invoices

**`GET /api/v1/invoices`**

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 10)
- `search` - Filter by service name
- `sort_by` - Field to sort by
- `sort_dir` - Sort direction (`asc`/`desc`)

#### Get Invoice

**`GET /api/v1/invoices/{id}`**

#### Create Invoice

**`POST /api/v1/invoices`**

**Request Body:**
```json
{
  "service_name": "DMP Service",
  "invoice_number": 1001,
  "date": "2024-03-16T00:00:00Z",
  "amount": 1500.50,
  "status": "Pending"
}
```

#### Update Invoice

**`PUT /api/v1/invoices/{id}`**

**Request Body:** Same as Create.

#### Delete Invoice

**`DELETE /api/v1/invoices/{id}`**

---

## 🧪 Example Usage

### Step 1: Create an Invoice
```bash
curl -X POST http://localhost:3000/api/v1/invoices -H "Content-Type: application/json" -d '{
    "service_name": "DMP Service",
    "invoice_number": 1001,
    "date": "2024-03-16T00:00:00Z",
    "amount": 1500.50,
    "status": "Pending"
}'
```

### Step 2: List Invoices
```bash
curl http://localhost:3000/api/v1/invoices?page=1&limit=10
```

### Step 3: Get Invoice Details
```bash
curl http://localhost:3000/api/v1/invoices/1
```

---

## 🔴 Error Responses

The API uses standard HTTP status codes and returns errors in the following format:
```json
{
  "message": "Error description",
  "error": "Detailed error message"
}
```

---

## 📊 Monitoring & Metrics

### Prometheus Metrics
Access metrics at: `http://localhost:3000/metrics`

Available metrics:
- `http_requests_total` - Total number of HTTP requests
- `http_request_duration_seconds` - HTTP request duration in seconds

### Grafana Dashboard
- URL: http://localhost:3001
- Login: admin/admin
- Default dashboard includes:
  - Request rates
  - Response times
  - HTTP status codes

### Prometheus UI
- URL: http://localhost:9090
- Query examples:
  - `rate(http_requests_total[1m])`
  - `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))`

---

## 🧩 Development Features

- Health check endpoint
- Input validation
- Error handling middleware
- Request logging
- Panic recovery
- CORS support
- Graceful shutdown
- Database connection retry logic
- Swagger documentation

