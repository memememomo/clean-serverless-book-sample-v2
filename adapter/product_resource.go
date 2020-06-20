package adapter

import (
	"clean-serverless-book-sample/domain"
	"time"
)

// ProductResource ProductModelのDynamoDB用構造体
type ProductResource struct {
	ResourceSchema
	DynamoResourceBase
	domain.ProductModel
	Mapper *DynamoModelMapper `dynamo:"-"`
}

func NewProductResource(productModel *domain.ProductModel, mapper *DynamoModelMapper) *ProductResource {
	return &ProductResource{
		ProductModel: *productModel,
		Mapper:       mapper,
	}
}

// 以下、DynamoResourceインタフェースの実装

// EntityName エンティティ名を返す。構造体名をエンティティ名として返すように実装している
func (p *ProductResource) EntityName() string {
	return p.Mapper.GetEntityNameFromStruct(*p)
}

// PK DynamoDBのHASHキーとして指定する文字列を返す
func (p *ProductResource) PK() string {
	return p.Mapper.GetPK(p)
}

// SetPK DynamoDBのHASHキーを設定する
func (p *ProductResource) SetPK() {
	p.ResourceSchema.PK = p.PK()
}

// SK DynamoDBのRANGEキーとして指定する文字列を返す
func (p *ProductResource) SK() string {
	return p.Mapper.GetSK(p)
}

// SetSK DynamoDBのRANGEキーを設定する
func (p *ProductResource) SetSK() {
	p.ResourceSchema.SK = p.SK()
}

// ID エンティティID
func (p *ProductResource) ID() uint64 {
	return p.ProductModel.ID
}

// SetID エンティティIDを設定する
func (p *ProductResource) SetID(id uint64) {
	p.ProductModel.ID = id
}

// Version DynamoDBの作成/更新をするときに楽観的ロックを行うためのVersionを返す
func (p *ProductResource) Version() int {
	return p.DynamoResourceBase.Version
}

// SetVersion DynamoDBの作成?更新をするときに楽観的ロックを行うためのVersionを設定する
func (p *ProductResource) SetVersion(v int) {
	p.DynamoResourceBase.Version = v
}

// CreatedAt レコードの作成時刻を返す
func (p *ProductResource) CreatedAt() time.Time {
	return p.DynamoResourceBase.CreatedAt
}

// SetCreatedAt レコードの作成時刻を設定する
func (p *ProductResource) SetCreatedAt(t time.Time) {
	p.DynamoResourceBase.CreatedAt = t
}

// UpdatedAt レコードの更新時刻を返す
func (p *ProductResource) UpdatedAt() time.Time {
	return p.DynamoResourceBase.UpdatedAt
}

// SetUpdatedAt レコードの更新時刻を設定する
func (p *ProductResource) SetUpdatedAt(t time.Time) {
	p.DynamoResourceBase.UpdatedAt = t
}
