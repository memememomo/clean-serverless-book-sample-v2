package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

// DeleteMicropost マイクロポスト削除
type DeleteMicropost struct {
	Getter              usecase.IGetMicropostByID
	MicropostRepository domain.MicropostRepository
}

func NewDeleteMicropost(getter usecase.IGetMicropostByID, repos domain.MicropostRepository) *DeleteMicropost {
	return &DeleteMicropost{
		Getter:              getter,
		MicropostRepository: repos,
	}
}

// Execute マイクロポストを削除
func (m *DeleteMicropost) Execute(req *usecase.DeleteMicropostRequest) (*usecase.DeleteMicropostResponse, error) {
	res, err := m.Getter.Execute(&usecase.GetMicropostByIDRequest{
		MicropostID: req.MicropostID,
		UserID:      req.UserID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = m.MicropostRepository.DeleteMicropost(res.Micropost.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.DeleteMicropostResponse{}, nil
}
