package main

import (
	"fmt"
	"net/http"

	"source-asia/internal/api"
	"source-asia/internal/products"
	"source-asia/internal/ratelimit"
)

func main() {

	limiter := ratelimit.NewLimiter()

	store := products.NewStore()

	http.HandleFunc("/request", api.RequestHandler(limiter))

	http.HandleFunc("/stats", api.StatsHandler(limiter))

	http.HandleFunc("/products", products.ProductHandler(store))

	http.HandleFunc("/products/", products.ProductDetailHandler(store))

	fmt.Println("Server listening on :8080")

	http.ListenAndServe(":8080", nil)
} 

