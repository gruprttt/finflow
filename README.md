#  FinFlow — Distributed System with Observability (SRE Project)

##  Project Overview

FinFlow is a **microservices-based distributed system** built to simulate real-world backend architecture with **observability, caching, messaging, and monitoring**.

This project demonstrates **SRE concepts** like:

* Metrics (Prometheus)
* Logging (Loki + Promtail)
* Dashboards (Grafana)
* Caching (Redis)
* Messaging (Kafka)
* Idempotency handling
* Latency & throughput monitoring

---

##  Architecture

```
User → Node API → Go Service → MySQL
                        ↓
                      Redis (cache)
                        ↓
                      Kafka → Python Consumer
```

### Observability Layer:

```
Prometheus → Metrics
Grafana → Dashboards
Loki + Promtail → Logs
```

---

## 🛠️ Tech Stack

| Layer            | Technology      |
| ---------------- | --------------- |
| API Layer        | Node.js         |
| Core Service     | Go (Golang)     |
| Database         | MySQL           |
| Cache            | Redis           |
| Messaging        | Kafka           |
| Consumer         | Python          |
| Metrics          | Prometheus      |
| Logs             | Loki + Promtail |
| Visualization    | Grafana         |
| Containerization | Docker Compose  |

---

##  Features Implemented

###  Backend Features

* User creation API
* Order creation API
* Idempotency handling (duplicate prevention)
* Redis caching for read optimization
* Kafka event publishing (`order_created`)

###  Observability

* Prometheus metrics:

  * `http_requests_total`
  * `http_request_duration_seconds`
  * `orders_created_total`
* Grafana dashboards
* Loki log aggregation
* Promtail log collection

###  Performance Monitoring

* Request latency tracking (Histogram)
* Throughput tracking (Requests/sec)
* Error rate monitoring (5xx tracking)

---

##  Project Structure

```
finflow/
│
├── services/
│   ├── node-api/
│   ├── go-service/
│   └── python-consumer/
│
├── infra/
│   └── docker-compose/
│       ├── docker-compose.yaml
│       ├── prometheus.yml
│       ├── promtail-config.yml
│       └── mysql-init/
│
└── README.md
```

---

##  Getting Started

### 1️ Clone Repository

```bash
git clone <your-repo-url>
cd finflow/infra/docker-compose
```

---

### 2️ Start System

```bash
docker compose up --build -d
```

---

### 3️ Verify Services

```bash
docker ps
```

---

##  API Testing

###  Create User

```bash
curl -X POST http://localhost:3000/users \
-H "Content-Type: application/json" \
-d '{"name":"Gurpreet","email":"test@mail.com"}'
```

---

### 🔹 Create Order

```bash
curl -X POST http://localhost:3000/orders \
-H "Content-Type: application/json" \
-d '{"user_id":1,"amount":5000,"idempotency_key":"order-1"}'
```

---

### 🔹 Idempotency Check

```bash
# Same request again
→ "Duplicate request ignored"
```

---

### 🔹 Get Order

```bash
curl http://localhost:3000/orders/1
```

---

##  Monitoring

### 🔹 Prometheus

```
http://localhost:9090
```

### 🔹 Grafana

```
http://localhost:3001
Login: admin / admin
```

### 🔹 Loki Logs

Available via Grafana → Explore → Loki

---

## 📈 Key Metrics (PromQL)

### Throughput

```
rate(http_requests_total[1m])
```

### Error Rate

```
rate(http_requests_total{status=~"5.."}[1m])
```

### Latency (P95)

```
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

---

## 🧠 Key Concepts Demonstrated

* Microservices architecture
* Event-driven system (Kafka)
* Caching strategy (Redis)
* Idempotent APIs
* Observability (metrics + logs)
* Latency & performance monitoring
* Containerized deployment

---

##  Known Limitations (Planned Improvements)

* DLQ (Dead Letter Queue) not implemented yet
* No retry mechanism for Kafka failures
* No alerting system yet
* No authentication layer
* Not deployed on Kubernetes / Cloud

---

##  Future Enhancements

* [ ] Implement DLQ (Kafka failure handling)
* [ ] Add retry with exponential backoff
* [ ] Grafana alerting (error rate, latency)
* [ ] Structured logging with correlation IDs
* [ ] Load testing (k6 / wrk)
* [ ] Kubernetes deployment
* [ ] Cloud deployment (AWS/GCP)

---

##  Learning Outcomes

This project helped in understanding:

* How real distributed systems work
* Observability in production systems
* Handling failures in async systems
* Monitoring latency, throughput, and errors
* End-to-end system design
