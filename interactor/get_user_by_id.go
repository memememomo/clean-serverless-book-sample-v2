package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

type GetUserByID struct {
	UserRepository domain.UserRepository
}

func NewGetUserByID(repos domain.UserRepository) *GetUserByID {
	return &GetUserByID{
		UserRepository: repos,
	}
}

// Execute ユーザーを取得
func (u *GetUserByID) Execute(req *usecase.GetUserByIDRequest) (*usecase.GetUserByIDResponse, error) {
	user, err := u.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.GetUserByIDResponse{User: user}, nil
}
