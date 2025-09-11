package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone     string `json:"phone" gorm:"uniqueIndex"`
	Code      int
	SessionId string `json:"sessionId"`
	Orders    []Order
}

func NewUser(phone string) *User {
	user := &User{
		Phone: phone,
		Code:  3245,
	}
	user.GenerateSessionId()
	return user
}

func (user *User) GenerateSessionId() {
	user.SessionId = uuid.NewString()
}

type Product struct {
	gorm.Model
	Name        string
	Description string
	// Images      pq.StringArray
}

func NewProduct(name, description string) *Product {
	return &Product{
		Name:        name,
		Description: description,
	}
}

type Order struct {
	gorm.Model
	UserID   uint
	Products []Product `gorm:"many2many:product_orders;"`
}

type ProductOrder struct {
	ProductID uint
	OrderID   uint
}
