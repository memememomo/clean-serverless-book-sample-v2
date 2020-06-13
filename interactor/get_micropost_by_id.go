package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

type GetMicropostByID struct {
	MicropostRepository domain.MicropostRepository
}

func NewGetMicropostByID(repos domain.MicropostRepository) *GetMicropostByID {
	return &GetMicropostByID{MicropostRepository: repos}
}

// GetMicropostByID マイクロポスト取得
func (m *GetMicropostByID) Execute(req *usecase.GetMicropostByIDRequest) (*usecase.GetMicropostByIDResponse, error) {
	micropost, err := m.MicropostRepository.GetMicropostByID(req.MicropostID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if micropost.UserID != req.UserID {
		return nil, errors.WithStack(domain.ErrNotFound)
	}
	return &usecase.GetMicropostByIDResponse{Micropost: micropost}, nil
}
