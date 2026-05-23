package products

type Product struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	SKU        string   `json:"sku"`
	ImageURLs  []string `json:"image_urls"`
	VideoURLs  []string `json:"video_urls"`
}

type ProductListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	ImageCount  int    `json:"image_count"`
	VideoCount  int    `json:"video_count"`
} 
