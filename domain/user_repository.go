package domain

// UserRepository ユーザーモデルのリポジトリ
type UserRepository interface {
	GetUsers() ([]*UserModel, error)
	GetUserByID(id uint64) (*UserModel, error)
	GetUserByEmail(email string) (*UserModel, error)
	CreateUser(newUser *UserModel) (*UserModel, error)
	UpdateUser(newUser *UserModel) error
	DeleteUser(targetUser *UserModel) error
}
