package usecase

// IDeleteUser ユーザー削除UseCase
type IDeleteUser interface {
	Execute(req *DeleteUserRequest) (*DeleteUserResponse, error)
}

// DeleteUserRequest ユーザー削除Request
type DeleteUserRequest struct {
	UserID uint64
}

// DeleteUserResponse ユーザー削除Response
type DeleteUserResponse struct {
}
