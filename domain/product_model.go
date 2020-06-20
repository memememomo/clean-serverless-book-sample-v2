package domain

import "time"

// ProductModel 製品を表すModel
type ProductModel struct {
	ID          uint64
	Name        string
	Price       int
	ReleaseDate time.Time
}

func NewProductModel(name string, price int, releaseDate time.Time) *ProductModel {
	return &ProductModel{
		Name:        name,
		Price:       price,
		ReleaseDate: releaseDate,
	}
}
