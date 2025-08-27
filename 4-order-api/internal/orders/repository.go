package orders

import (
	"demo/order/4-order-api/pkg/db"

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

func (repo *ProductRepository) AddProduct(product *Product) (*Product, error) {
	err := repo.Database.DB.Create(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *ProductRepository) UpdateProduct(product *Product) (*Product, error) {
	err := repo.Database.DB.Clauses(clause.Returning{}).Updates(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *ProductRepository) DeleteProductById(id uint) error {
	err := repo.Database.DB.Delete(&Product{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) GetProductById(id uint) (*Product, error) {
	var product Product
	err := repo.Database.DB.First(&product, id).Error
	return &product, err
}
