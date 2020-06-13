package usecase

type ICreateMicropost interface {
	Execute(req *CreateMicropostRequest) (*CreateMicropostResponse, error)
}

type CreateMicropostRequest struct {
	Content string
	UserID  uint64
}

type CreateMicropostResponse struct {
	MicropostID uint64
}
