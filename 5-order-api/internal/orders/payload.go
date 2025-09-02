package orders

type ProductCreateRequest struct {
	Name        string `json:"name"  validate:"required"`
	Description string `json:"description"  validate:"required"`
}

type ProductUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
