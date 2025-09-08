package orders

import (
	"demo/order/6-order-api-cart/pkg/db"

	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (repo *ProductRepository) AddProduct(product *db.Product) (*db.Product, error) {
	err := repo.Database.DB.Create(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *ProductRepository) UpdateProduct(product *db.Product) (*db.Product, error) {
	err := repo.Database.DB.Clauses(clause.Returning{}).Updates(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *ProductRepository) DeleteProductById(id uint) error {
	err := repo.Database.DB.Delete(&db.Product{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) GetProductById(id uint) (*db.Product, error) {
	var product db.Product
	err := repo.Database.DB.First(&product, id).Error
	return &product, err
}

func (repo *ProductRepository) GetProductByName(name string) (*db.Product, error) {
	var product db.Product
	err := repo.Database.DB.Where(&db.Product{Name: name}).First(&product).Error
	return &product, err
}

func (repo *ProductRepository) AddOrder(userId uint) (*db.Order, error) {
	order := &db.Order{
		UserID: userId,
	}
	err := repo.Database.DB.Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *ProductRepository) AddProductToOrder(order *db.Order, product *db.Product) error {
	err := repo.Database.DB.Model(&order).Association("Products").Append(&product)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) GetOrdersByUserId(userId uint) (*[]db.Order, error) {
	var orders []db.Order
	err := repo.Database.DB.Where(&db.Order{UserID: userId}).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func (repo *ProductRepository) GetOrderByOrderId(userId, orderId uint) (*db.Order, error) {
	var order db.Order
	err := repo.Database.DB.Preload("Products").Where(db.Order{UserID: userId}).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
