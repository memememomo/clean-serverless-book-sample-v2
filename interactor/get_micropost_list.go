package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

// GetMicropostList マイクロポスト取得
type GetMicropostList struct {
	MicropostRepository domain.MicropostRepository
}

func NewGetMicropostList(repos domain.MicropostRepository) *GetMicropostList {
	return &GetMicropostList{
		MicropostRepository: repos,
	}
}

// Execute マイクロポスト一覧取得
func (m *GetMicropostList) Execute(req *usecase.GetMicropostListRequest) (*usecase.GetMicropostListResponse, error) {
	microposts, err := m.MicropostRepository.GetMicropostsByUserID(req.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.GetMicropostListResponse{Microposts: microposts}, nil
}
