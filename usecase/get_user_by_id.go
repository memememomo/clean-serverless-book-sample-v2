package usecase

import "clean-serverless-book-sample-v2/domain"

// IGetUserByID 指定されたIDのユーザーを取得UseCase
type IGetUserByID interface {
	Execute(req *GetUserByIDRequest) (*GetUserByIDResponse, error)
}

type GetUserByIDRequest struct {
	UserID uint64
}

type GetUserByIDResponse struct {
	User *domain.UserModel
}
