package adapter

import (
	"clean-serverless-book-sample-v2/domain"
	"time"
)

// UserResource DynamoDB上のデータ構造を表した構造体
type UserResource struct {
	ResourceSchema
	DynamoResourceBase
	domain.UserModel
	Mapper *DynamoModelMapper `dynamo:"-"`
}

func NewUserResource(userModel *domain.UserModel, mapper *DynamoModelMapper) *UserResource {
	return &UserResource{
		UserModel: *userModel,
		Mapper:    mapper,
	}
}

// DynamoResourceインタフェースの実装

func (u *UserResource) EntityName() string {
	return u.Mapper.GetEntityNameFromStruct(*u)
}

func (u *UserResource) PK() string {
	return u.Mapper.GetPK(u)
}

func (u *UserResource) SetPK() {
	u.ResourceSchema.PK = u.PK()
}

func (u *UserResource) SK() string {
	return u.Mapper.GetSK(u)
}

func (u *UserResource) SetSK() {
	u.ResourceSchema.SK = u.SK()
}

func (u *UserResource) SetID(id uint64) {
	u.UserModel.ID = id
}

func (u *UserResource) ID() uint64 {
	return u.UserModel.ID
}

func (u *UserResource) SetVersion(v int) {
	u.DynamoResourceBase.Version = v
}

func (u *UserResource) Version() int {
	return u.DynamoResourceBase.Version
}

func (u *UserResource) CreatedAt() time.Time {
	return u.DynamoResourceBase.CreatedAt
}

func (u *UserResource) SetCreatedAt(t time.Time) {
	u.DynamoResourceBase.CreatedAt = t
}

func (u *UserResource) UpdatedAt() time.Time {
	return u.DynamoResourceBase.UpdatedAt
}

func (u *UserResource) SetUpdatedAt(t time.Time) {
	u.DynamoResourceBase.UpdatedAt = t
}
