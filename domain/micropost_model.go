package domain

// MicropostModel マイクロポストのモデル
type MicropostModel struct {
	ID      uint64
	Content string
	UserID  uint64
}

func NewMicropostModel(content string, userID uint64) *MicropostModel {
	return &MicropostModel{Content: content, UserID: userID}
}
