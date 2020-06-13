package adapter

import (
	"clean-serverless-book-sample-v2/domain"
	"github.com/guregu/dynamo"
	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
)

// MicropostOperator マイクロポストを操作する構造体
type MicropostOperator struct {
	Client *ResourceTableOperator
	Mapper *DynamoModelMapper
}

func (m *MicropostOperator) getMicropostResourceByID(id uint64) (*MicropostResource, error) {
	var micropostResource MicropostResource
	_, err := m.Mapper.GetEntityByID(id, &MicropostResource{}, &micropostResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &micropostResource, nil
}

// GetMicropostByID IDでマイクロポストを取得する
func (m *MicropostOperator) GetMicropostByID(id uint64) (*domain.MicropostModel, error) {
	micropostResource, err := m.getMicropostResourceByID(id)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return nil, errors.WithStack(domain.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &micropostResource.MicropostModel, nil
}

// GetMicropostsByUserID 指定されたユーザーIDに紐づいているマイクロポスト一覧を取得する
func (m *MicropostOperator) GetMicropostsByUserID(userID uint64) ([]*domain.MicropostModel, error) {
	table, err := m.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fb := nomof.NewBuilder()
	fb.Equal("UserID", userID)
	fb.BeginsWith("PK", m.Mapper.GetEntityNameFromStruct(MicropostResource{}))

	var micropostResource []MicropostResource
	err = table.Scan().Filter(fb.JoinAnd(), fb.Arg...).All(&micropostResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var microposts = make([]*domain.MicropostModel, len(micropostResource))
	for i := range micropostResource {
		microposts[i] = &micropostResource[i].MicropostModel
	}

	return microposts, nil
}

// DeleteMicropost 指定されたIDのマイクロポストを削除する
func (m *MicropostOperator) DeleteMicropost(id uint64) error {
	micropost, err := m.getMicropostResourceByID(id)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return errors.WithStack(domain.ErrNotFound)
		}
		return errors.WithStack(err)
	}

	err = m.Mapper.DeleteResource(micropost)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CreateMicropost 新規作成する
func (m *MicropostOperator) CreateMicropost(micropostModel *domain.MicropostModel) (*domain.MicropostModel, error) {
	micropostResource := NewMicropostResource(micropostModel, m.Mapper)
	err := m.Mapper.PutResource(micropostResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &micropostResource.MicropostModel, nil
}

// UpdateMicropost 更新する
func (m *MicropostOperator) UpdateMicropost(micropostModel *domain.MicropostModel) error {
	micropostResource, err := m.getMicropostResourceByID(micropostModel.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	micropostResource.Content = micropostModel.Content

	err = m.Mapper.PutResource(micropostResource)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
