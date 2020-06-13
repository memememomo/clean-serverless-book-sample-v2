package adapter

import (
	"clean-serverless-book-sample-v2/domain"
	"time"
)

// MicropostResource DynamoDB上のデータ構造を表した構造体
type MicropostResource struct {
	ResourceSchema
	DynamoResourceBase
	domain.MicropostModel
	Mapper *DynamoModelMapper `dynamo:"-"`
}

func NewMicropostResource(micropostModel *domain.MicropostModel, mapper *DynamoModelMapper) *MicropostResource {
	return &MicropostResource{
		MicropostModel: *micropostModel,
		Mapper:         mapper,
	}
}

// DynamoResourceインタフェースの実装

func (m *MicropostResource) EntityName() string {
	return m.Mapper.GetEntityNameFromStruct(*m)
}

func (m *MicropostResource) PK() string {
	return m.Mapper.GetPK(m)
}

func (m *MicropostResource) SetPK() {
	m.ResourceSchema.PK = m.PK()
}

func (m *MicropostResource) SK() string {
	return m.Mapper.GetSK(m)
}

func (m *MicropostResource) SetSK() {
	m.ResourceSchema.SK = m.SK()
}

func (m *MicropostResource) SetID(id uint64) {
	m.MicropostModel.ID = id
}

func (m *MicropostResource) ID() uint64 {
	return m.MicropostModel.ID
}

func (m *MicropostResource) SetVersion(v int) {
	m.DynamoResourceBase.Version = v
}

func (m *MicropostResource) Version() int {
	return m.DynamoResourceBase.Version
}

func (m *MicropostResource) CreatedAt() time.Time {
	return m.DynamoResourceBase.CreatedAt
}

func (m *MicropostResource) SetCreatedAt(t time.Time) {
	m.DynamoResourceBase.CreatedAt = t
}

func (m *MicropostResource) UpdatedAt() time.Time {
	return m.DynamoResourceBase.UpdatedAt
}

func (m *MicropostResource) SetUpdatedAt(t time.Time) {
	m.DynamoResourceBase.UpdatedAt = t
}
