package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

// GetUserList ユーザー取得
type GetUserList struct {
	UserRepository domain.UserRepository
}

func NewGetUserList(repos domain.UserRepository) *GetUserList {
	return &GetUserList{
		UserRepository: repos,
	}
}

// Execute ユーザー一覧を取得
func (u *GetUserList) Execute(req *usecase.GetUserListRequest) (*usecase.GetUserListResponse, error) {
	users, err := u.UserRepository.GetUsers()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.GetUserListResponse{Users: users}, nil
}
