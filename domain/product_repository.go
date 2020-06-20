package domain

type ProductRepository interface {
	CreateProduct(newProduct *ProductModel) (*ProductModel, error)
	UpdateProduct(newProduct *ProductModel) error
	GetProductByID(id uint64) (*ProductModel, error)
	GetProducts() ([]*ProductModel, error)
	DeleteProduct(id uint64) error
}
