package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

// UpdateMicropost
type UpdateMicropost struct {
	MicropostRepository domain.MicropostRepository
}

func NewUpdateMicropost(repos domain.MicropostRepository) *UpdateMicropost {
	return &UpdateMicropost{
		MicropostRepository: repos,
	}
}

// Execute 更新
func (m *UpdateMicropost) Execute(req *usecase.UpdateMicropostRequest) (*usecase.UpdateMicropostResponse, error) {
	newMicropost := domain.NewMicropostModel(req.Content, req.UserID)
	newMicropost.ID = req.MicropostID
	err := m.MicropostRepository.UpdateMicropost(newMicropost)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.UpdateMicropostResponse{}, nil
}
