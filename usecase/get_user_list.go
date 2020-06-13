package usecase

import "clean-serverless-book-sample-v2/domain"

// IGetUserList ユーザー一覧取得UseCase
type IGetUserList interface {
	Execute(req *GetUserListRequest) (*GetUserListResponse, error)
}

// GetUserListRequest ユーザー一覧取得Request
type GetUserListRequest struct {
}

// GetUserListResponse ユーザー一覧取得Response
type GetUserListResponse struct {
	Users []*domain.UserModel
}

func (g *GetUserListResponse) UserCount() int {
	return len(g.Users)
}
