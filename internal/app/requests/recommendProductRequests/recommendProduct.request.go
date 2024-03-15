package recommendProductRequests

type RecommendProductResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}
