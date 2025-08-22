package orders

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string
}
