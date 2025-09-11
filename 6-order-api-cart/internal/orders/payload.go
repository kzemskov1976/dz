package orders

type ProductCreateRequest struct {
	Name        string `json:"name"  validate:"required"`
	Description string `json:"description"  validate:"required"`
}

type ProductUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OrderCreateRequest struct {
	Products []string `json:"products" validate:"required"`
}

type OrderCreateResponse struct {
	OrderId uint `json:"orderId"`
}

type MyOrdersResponse struct {
	OrderId []uint `json:"orders"`
}

type OrderResponse struct {
	OrderId  uint     `json:"orderId"`
	Products []string `json:"products"`
}
