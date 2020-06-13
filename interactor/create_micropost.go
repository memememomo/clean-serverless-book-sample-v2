package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

// CreateMicropost マイクロポスト作成
type CreateMicropost struct {
	MicropostRepository domain.MicropostRepository
}

func NewCreateMicropost(repos domain.MicropostRepository) *CreateMicropost {
	return &CreateMicropost{
		MicropostRepository: repos,
	}
}

// Execute マイクロポストを新規作成
func (m *CreateMicropost) Execute(req *usecase.CreateMicropostRequest) (*usecase.CreateMicropostResponse, error) {
	newMicropost := domain.NewMicropostModel(req.Content, req.UserID)
	micropost, err := m.MicropostRepository.CreateMicropost(newMicropost)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.CreateMicropostResponse{MicropostID: micropost.ID}, nil
}
