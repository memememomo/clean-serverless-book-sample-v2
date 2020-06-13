package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

// UpdateUser ユーザー更新
type UpdateUser struct {
	UserRepository domain.UserRepository
	UniqChecker    *domain.UserEmailUniqChecker
}

func NewUpdateUser(repos domain.UserRepository, checker *domain.UserEmailUniqChecker) *UpdateUser {
	return &UpdateUser{
		UserRepository: repos,
		UniqChecker:    checker,
	}
}

// Execute ユーザーを更新
func (u *UpdateUser) Execute(req *usecase.UpdateUserRequest) (*usecase.UpdateUserResponse, error) {
	isUniq, err := u.UniqChecker.IsUniqueEmail(req.ToUserModel())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !isUniq {
		return nil, errors.WithStack(ErrUniqEmail)
	}

	err = u.UserRepository.UpdateUser(req.ToUserModel())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.UpdateUserResponse{}, nil
}
