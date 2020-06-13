package interactor

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/pkg/errors"
)

// UserDeleter ユーザー削除
type UserDeleter struct {
	UserRepository domain.UserRepository
	UserGetter     usecase.IGetUserByID
}

func NewUserDeleter(repos domain.UserRepository, getter usecase.IGetUserByID) *UserDeleter {
	return &UserDeleter{
		UserRepository: repos,
		UserGetter:     getter,
	}
}

// Execute ユーザーを削除
func (u *UserDeleter) Execute(req *usecase.DeleteUserRequest) (*usecase.DeleteUserResponse, error) {
	user, err := u.UserGetter.Execute(&usecase.GetUserByIDRequest{UserID: req.UserID})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = u.UserRepository.DeleteUser(user.User)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.DeleteUserResponse{}, nil
}
