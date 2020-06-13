package usecase

type IUpdateMicropost interface {
	Execute(req *UpdateMicropostRequest) (*UpdateMicropostResponse, error)
}

type UpdateMicropostRequest struct {
	Content     string
	UserID      uint64
	MicropostID uint64
}

type UpdateMicropostResponse struct {
}
