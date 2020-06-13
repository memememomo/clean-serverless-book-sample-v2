package usecase

import "clean-serverless-book-sample-v2/domain"

type IUpdateUser interface {
	Execute(req *UpdateUserRequest) (*UpdateUserResponse, error)
}

type UpdateUserRequest struct {
	ID    uint64
	Name  string
	Email string
}

func (u *UpdateUserRequest) ToUserModel() *domain.UserModel {
	return &domain.UserModel{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

type UpdateUserResponse struct {
}
