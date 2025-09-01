package orders

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string
	Description string
	Images      pq.StringArray
}

func NewProduct(req ProductCreateRequest) *Product {
	return &Product{
		Name:        req.Name,
		Description: req.Description,
	}
}
