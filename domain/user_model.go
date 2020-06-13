package domain

// UserModel ユーザーモデル
type UserModel struct {
	ID    uint64
	Name  string
	Email string
}

func NewUserModel(name, email string) *UserModel {
	return &UserModel{Name: name, Email: email}
}
