package products

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func isValidURL(value string) bool {
	if len(value) > 2048 {
		return false
	}

	u, err := url.Parse(value)

	if err != nil {
		return false
	}

	return u.Scheme == "http" || u.Scheme == "https"
}

func validateURLs(urls []string) bool {
	if len(urls) > 20 {
		return false
	}

	for _, u := range urls {
		if !isValidURL(u) {
			return false
		}
	}

	return true
}

func ProductHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {

			var product Product

			err := json.NewDecoder(r.Body).Decode(&product)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid JSON",
				})
				return
			}

			if product.Name == "" || product.SKU == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "name and sku are required",
				})
				return
			}

			if !validateURLs(product.ImageURLs) || !validateURLs(product.VideoURLs) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid URLs",
				})
				return
			}

			created, err := store.CreateProduct(product)

			if err != nil {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(created)

			return
		}

		if r.Method == http.MethodGet {

			limit := 10
			offset := 0

			limitParam := r.URL.Query().Get("limit")
			offsetParam := r.URL.Query().Get("offset")

			if limitParam != "" {
				value, err := strconv.Atoi(limitParam)
				if err == nil {
					limit = value
				}
			}

			if offsetParam != "" {
				value, err := strconv.Atoi(offsetParam)
				if err == nil {
					offset = value
				}
			}

			items := store.ListProducts(limit, offset)

			json.NewEncoder(w).Encode(items)

			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ProductDetailHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/products/")

		if strings.HasSuffix(path, "/media") {

			id := strings.TrimSuffix(path, "/media")

			var body struct {
				ImageURLs []string `json:"image_urls"`
				VideoURLs []string `json:"video_urls"`
			}

			err := json.NewDecoder(r.Body).Decode(&body)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if len(body.ImageURLs) == 0 && len(body.VideoURLs) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !validateURLs(body.ImageURLs) || !validateURLs(body.VideoURLs) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			product, err := store.AddMedia(id, body.ImageURLs, body.VideoURLs)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(product)

			return
		}

		id := path

		product, exists := store.GetProduct(id)

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(product)
	}
} 
