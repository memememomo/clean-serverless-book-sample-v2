package usecase

import "clean-serverless-book-sample-v2/domain"

// ICreateUser ユーザー新規作成UseCase
type ICreateUser interface {
	Execute(req *CreateUserRequest) (*CreateUserResponse, error)
}

// CreateUserRequest ユーザー新規作成Request
type CreateUserRequest struct {
	Name  string
	Email string
}

func (u *CreateUserRequest) ToUserModel() *domain.UserModel {
	return domain.NewUserModel(u.Name, u.Email)
}

// CreateUserResponse ユーザー新規作成Response
type CreateUserResponse struct {
	User *domain.UserModel
}

func (u *CreateUserResponse) GetUserID() uint64 {
	return u.User.ID
}
