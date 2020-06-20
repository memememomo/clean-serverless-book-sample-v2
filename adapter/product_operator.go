package adapter

import (
	"clean-serverless-book-sample-v2/domain"
	"github.com/guregu/dynamo"
	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
)

// ProductOperator ProductModelのCRUDを実装する構造体
type ProductOperator struct {
	Client *ResourceTableOperator
	Mapper *DynamoModelMapper
}

func (p *ProductOperator) getProductResourceByID(id uint64) (*ProductResource, error) {
	var productResource ProductResource
	_, err := p.Mapper.GetEntityByID(id, &ProductResource{}, &productResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &productResource, nil
}

// GetProductByID IDによるProduct取得処理
func (p *ProductOperator) GetProductByID(id uint64) (*domain.ProductModel, error) {
	// IDによるProduct取得処理
	productResource, err := p.getProductResourceByID(id)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return nil, errors.WithStack(domain.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	// Productを返す
	return &productResource.ProductModel, nil
}

// GetProducts 一覧取得処理
func (p *ProductOperator) GetProducts() ([]*domain.ProductModel, error) {
	// DynamoDBテーブルに接続するためのクライアントを取得
	table, err := p.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// フィルタの設定
	fb := nomof.NewBuilder()
	fb.BeginsWith("PK", p.Mapper.GetEntityNameFromStruct(ProductResource{}))

	// DynamoDBから一覧取得処理
	var productResource []ProductResource
	err = table.
		Scan().
		Filter(fb.JoinAnd(), fb.Arg...).
		All(&productResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// ProductResourceからProductModelに変換
	var products = make([]*domain.ProductModel, len(productResource))
	for i := range productResource {
		products[i] = &productResource[i].ProductModel
	}

	// 一覧を返す
	return products, nil
}

// CreateProduct 新規作成
func (p *ProductOperator) CreateProduct(productModel *domain.ProductModel) (*domain.ProductModel, error) {
	// ProductModelからProductResourceを作成する
	productResource := NewProductResource(productModel, p.Mapper)

	// DynamoDBに保存
	err := p.Mapper.PutResource(productResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 新規作成したProductModelを返す
	return &productResource.ProductModel, nil
}

// UpdateProduct 更新処理
func (p *ProductOperator) UpdateProduct(productModel *domain.ProductModel) error {
	// 既存のProductを取得する
	productResource, err := p.getProductResourceByID(productModel.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	// 更新内容をModelに反映
	productResource.ProductModel.Name = productModel.Name
	productResource.ProductModel.Price = productModel.Price
	productResource.ProductModel.ReleaseDate = productModel.ReleaseDate

	// 更新処理
	err = p.Mapper.PutResource(productResource)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// DeleteProduct 削除処理
func (p *ProductOperator) DeleteProduct(id uint64) error {
	// 既存のProductを取得する
	product, err := p.getProductResourceByID(id)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return errors.WithStack(domain.ErrNotFound)
		}
		return errors.WithStack(err)
	}

	// 削除処理
	err = p.Mapper.DeleteResource(product)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
