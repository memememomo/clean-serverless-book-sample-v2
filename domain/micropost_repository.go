package domain

// MicropostRepository Micropostモデルのリポジトリ
type MicropostRepository interface {
	CreateMicropost(newMicropost *MicropostModel) (*MicropostModel, error)
	UpdateMicropost(newMicropost *MicropostModel) error
	GetMicropostByID(id uint64) (*MicropostModel, error)
	GetMicropostsByUserID(userID uint64) ([]*MicropostModel, error)
	DeleteMicropost(id uint64) error
}
