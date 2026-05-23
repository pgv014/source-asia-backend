package products

import (
	"errors"
	"fmt"
	"sync"
)

type Store struct {
	mu sync.Mutex

	products map[string]Product
	skuIndex map[string]bool

	counter int
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]Product),
		skuIndex: make(map[string]bool),
		counter:  0,
	}
}

func (s *Store) CreateProduct(product Product) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.skuIndex[product.SKU] {
		return Product{}, errors.New("sku already exists")
	}

	s.counter++

	product.ID = fmt.Sprintf("prod_%06d", s.counter)

	s.products[product.ID] = product

	s.skuIndex[product.SKU] = true

	return product, nil
}

func (s *Store) GetProduct(id string) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[id]

	return product, exists
}

func (s *Store) ListProducts(limit int, offset int) []ProductListItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := []ProductListItem{}

	count := 0

	for _, product := range s.products {

		if count < offset {
			count++
			continue
		}

		if len(items) >= limit {
			break
		}

		item := ProductListItem{
			ID:         product.ID,
			Name:       product.Name,
			SKU:        product.SKU,
			ImageCount: len(product.ImageURLs),
			VideoCount: len(product.VideoURLs),
		}

		items = append(items, item)

		count++
	}

	return items
}

func (s *Store) AddMedia(id string, imageURLs []string, videoURLs []string) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[id]

	if !exists {
		return Product{}, errors.New("product not found")
	}

	product.ImageURLs = append(product.ImageURLs, imageURLs...)
	product.VideoURLs = append(product.VideoURLs, videoURLs...)

	s.products[id] = product

	return product, nil
} 

