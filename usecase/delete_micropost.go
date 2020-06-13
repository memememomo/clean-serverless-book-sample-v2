package usecase

type IDeleteMicropost interface {
	Execute(req *DeleteMicropostRequest) (*DeleteMicropostResponse, error)
}

type DeleteMicropostRequest struct {
	MicropostID uint64
	UserID      uint64
}

type DeleteMicropostResponse struct {
}
