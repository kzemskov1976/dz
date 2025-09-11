package di

import "demo/order/6-order-api-cart/pkg/db"

type IntUserRepository interface {
	FindByPhone(phone string) (*db.User, error)
}
