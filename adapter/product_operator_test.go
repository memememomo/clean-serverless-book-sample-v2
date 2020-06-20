package adapter_test

import (
	"clean-serverless-book-sample/adapter"
	"clean-serverless-book-sample/domain"
	"clean-serverless-book-sample/mocks"
	"clean-serverless-book-sample/registry"
	"fmt"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const DatetimeFormat = "2006-01-02T15:04:05Z"

func createProductResource(newProduct *domain.ProductModel) (*adapter.ProductResource, error) {
	client := registry.GetFactory().BuildResourceTableOperator()
	table, err := client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newProductResource := adapter.NewProductResource(
		newProduct,
		registry.GetFactory().BuildDynamoModelMapper())
	newProductResource.SetPK()
	newProductResource.SetSK()
	newProductResource.SetVersion(1)

	err = table.Put(newProductResource).Run()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return newProductResource, nil
}

func getProductResource(id uint64) (*adapter.ProductResource, error) {
	// DynamoDBにアクセスするためのクライアントを取得
	client := registry.GetFactory().BuildResourceTableOperator()

	// DynamoDBのテーブルにアクセスするためのクライアントを取得
	table, err := client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// DynamoDBテーブルから該当のレコードを取得
	var result adapter.ProductResource
	err = table.
		Get("PK", fmt.Sprintf("ProductResource-%011d", id)).
		Range("SK", dynamo.Equal, fmt.Sprintf("%011d", id)).
		One(&result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func TestProductOperator_CreateProduct(t *testing.T) {
	// テスト用のローカルDynamoDBを作成・接続
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 新規作成用のデータ
	product := domain.NewProductModel(
		"テスト", 100, time.Now())

	operator := registry.GetFactory().BuildProductOperator()

	// Createメソッドを呼び出し
	_, err := operator.CreateProduct(product)
	assert.NoError(t, err)

	// 作成されたレコードを取得できるか、内容は合っているかを確認
	result, err := getProductResource(1)
	assert.NoError(t, err)

	assert.Equal(t, product.Name, result.Name)
	assert.Equal(t, product.Price, result.Price)
	assert.Equal(t, product.ReleaseDate.Format(DatetimeFormat), result.ReleaseDate.Format(DatetimeFormat))
}

func TestProductOperator_UpdateProduct(t *testing.T) {
	// テスト用のローカルDynamoDBを作成・接続
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 更新用レコードを作成
	targetProduct, err := createProductResource(&domain.ProductModel{
		ID:          1,
		Name:        "テスト製品",
		Price:       100,
		ReleaseDate: time.Now(),
	})
	assert.NoError(t, err)

	// 更新パラメータ
	updatedProduct := &domain.ProductModel{
		ID:          targetProduct.ID(),
		Name:        "テスト製品(更新)",
		Price:       200,
		ReleaseDate: time.Now().AddDate(0, 0, 1),
	}

	// 更新処理
	operator := registry.GetFactory().BuildProductOperator()
	err = operator.UpdateProduct(updatedProduct)
	assert.NoError(t, err)

	// DynamoDBにあるデータが更新されているかチェック
	result, err := getProductResource(1)
	assert.NoError(t, err)

	assert.Equal(t, updatedProduct.Name, result.Name)
	assert.Equal(t, updatedProduct.Price, result.Price)
	assert.Equal(t, updatedProduct.ReleaseDate.Format(DatetimeFormat), result.ReleaseDate.Format(DatetimeFormat))
}

func TestProductOperator_GetProducts(t *testing.T) {
	// テスト用のローカルDynamoDBを作成・接続
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 取得用レコードを2つ作成
	product1, err := createProductResource(&domain.ProductModel{
		ID:          1,
		Name:        "テスト製品1",
		Price:       100,
		ReleaseDate: time.Now(),
	})
	assert.NoError(t, err)

	product2, err := createProductResource(&domain.ProductModel{
		ID:          2,
		Name:        "テスト製品2",
		Price:       200,
		ReleaseDate: time.Now().AddDate(0, 0, 1),
	})
	assert.NoError(t, err)

	// 一覧取得処理
	operator := registry.GetFactory().BuildProductOperator()
	products, err := operator.GetProducts()
	assert.NoError(t, err)

	// 所得した一覧の内容をチェック
	assert.Equal(t, product1.ID(), products[0].ID)
	assert.Equal(t, product1.Name, products[0].Name)
	assert.Equal(t, product1.ReleaseDate.Format(DatetimeFormat), products[0].ReleaseDate.Format(DatetimeFormat))

	assert.Equal(t, product2.ID(), products[1].ID)
	assert.Equal(t, product2.Name, products[1].Name)
	assert.Equal(t, product2.ReleaseDate.Format(DatetimeFormat), products[1].ReleaseDate.Format(DatetimeFormat))
}

func TestProductOperator_GetProductByID(t *testing.T) {
	// テスト用のローカルDynamoDBを作成・接続
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 取得用レコードを作成
	expected, err := createProductResource(&domain.ProductModel{
		ID:          1,
		Name:        "テスト製品",
		Price:       100,
		ReleaseDate: time.Now(),
	})

	// IDによるProduct取得処理
	operator := registry.GetFactory().BuildProductOperator()
	product, err := operator.GetProductByID(1)
	assert.NoError(t, err)

	// 取得した内容をチェック
	assert.Equal(t, expected.ID(), product.ID)
	assert.Equal(t, expected.Name, product.Name)
	assert.Equal(t, expected.ReleaseDate.Format(DatetimeFormat), product.ReleaseDate.Format(DatetimeFormat))
}

func TestProductOperator_DeleteProduct(t *testing.T) {
	// テスト用のローカルDynamoDBを作成・接続
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 削除用レコードを作成
	expected, err := createProductResource(&domain.ProductModel{
		ID:          1,
		Name:        "テスト製品",
		Price:       100,
		ReleaseDate: time.Now(),
	})
	assert.NoError(t, err)

	// 削除処理
	operator := registry.GetFactory().BuildProductOperator()
	err = operator.DeleteProduct(expected.ID())
	assert.NoError(t, err)

	// 削除されているかチェック
	_, err = getProductResource(1)
	assert.Equal(t, dynamo.ErrNotFound.Error(), err.Error())
}
