package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var rdb *redis.Client
var kafkaWriter *kafka.Writer
var ctx = context.Background()

// ------------------- METRICS -------------------

var ordersCreated = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "Total orders created",
	},
)

var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests",
	},
	[]string{"path", "method", "status"},
)

var httpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Request duration",
	},
	[]string{"path"},
)

// ------------------- INIT -------------------

func initDB() {
	var err error
	dsn := "root:root@tcp(mysql:3306)/finflow"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL ✅")
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	fmt.Println("Connected to Redis ✅")
}

func initKafka() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "order_created",
		Balancer: &kafka.LeastBytes{},
	}
	fmt.Println("Kafka Producer Ready ✅")
}

// ------------------- STRUCTS -------------------

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	UserID int    `json:"user_id"`
	Amount int    `json:"amount"`
	Key    string `json:"idempotency_key"`
}

// ------------------- HANDLERS -------------------

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Go Service 🚀"))
}

// Create User
func createUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	query := "INSERT INTO users(name, email) VALUES (?, ?)"
	_, err := db.Exec(query, user.Name, user.Email)
	if err != nil {
		httpRequestsTotal.WithLabelValues("/users", "POST", "500").Inc()
		http.Error(w, err.Error(), 500)
		return
	}

	httpRequestsTotal.WithLabelValues("/users", "POST", "200").Inc()

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("/users").Observe(duration)

	w.Write([]byte("User created"))
}

// Create Order (Idempotent + Kafka + Metrics)
func createOrder(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var order Order
	json.NewDecoder(r.Body).Decode(&order)

	if order.Key == "" {
		httpRequestsTotal.WithLabelValues("/orders", "POST", "400").Inc()
		http.Error(w, "idempotency_key required", 400)
		return
	}

	query := "INSERT INTO orders(user_id, amount, idempotency_key) VALUES (?, ?, ?)"
	_, err := db.Exec(query, order.UserID, order.Amount, order.Key)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			httpRequestsTotal.WithLabelValues("/orders", "POST", "200").Inc()
			w.Write([]byte("Duplicate request ignored"))
			return
		}

		httpRequestsTotal.WithLabelValues("/orders", "POST", "500").Inc()
		http.Error(w, err.Error(), 500)
		return
	}

	ordersCreated.Inc()
	httpRequestsTotal.WithLabelValues("/orders", "POST", "200").Inc()

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("/orders").Observe(duration)

	// Kafka event
	msg := fmt.Sprintf("UserID:%d Amount:%d", order.UserID, order.Amount)
	err = kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: []byte(msg),
	})

	if err != nil {
		fmt.Println("Kafka error:", err)
	}

	w.Write([]byte("Order created"))
}

// Get Order (Redis cache + DB)
func getOrder(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id := r.URL.Path[len("/orders/"):]

	// Check Redis
	val, err := rdb.Get(ctx, id).Result()
	if err == nil {
		httpRequestsTotal.WithLabelValues("/orders/:id", "GET", "200").Inc()
		w.Write([]byte("From Redis: " + val))
		return
	}

	// DB fallback
	query := "SELECT id, user_id, amount FROM orders WHERE id=?"
	row := db.QueryRow(query, id)

	var orderID, userID, amount int
	err = row.Scan(&orderID, &userID, &amount)
	if err != nil {
		httpRequestsTotal.WithLabelValues("/orders/:id", "GET", "500").Inc()
		http.Error(w, err.Error(), 500)
		return
	}

	result := fmt.Sprintf("OrderID:%d UserID:%d Amount:%d", orderID, userID, amount)

	// Save to Redis
	rdb.Set(ctx, id, result, 0)

	httpRequestsTotal.WithLabelValues("/orders/:id", "GET", "200").Inc()

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("/orders/:id").Observe(duration)

	w.Write([]byte("From DB: " + result))
}

// ------------------- MAIN -------------------

func main() {
	initDB()
	initRedis()
	initKafka()

	prometheus.MustRegister(ordersCreated)
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/users", createUser)
	http.HandleFunc("/orders", createOrder)
	http.HandleFunc("/orders/", getOrder)

	fmt.Println("Go service running on port 8080 🚀")
	http.ListenAndServe(":8080", nil)
}
