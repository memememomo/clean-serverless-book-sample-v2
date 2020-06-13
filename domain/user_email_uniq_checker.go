package domain

import (
	"github.com/pkg/errors"
)

// UserEmailUniqChecker メールアドレスの重複チェッカー
type UserEmailUniqChecker struct {
	Repos UserRepository
}

func NewUserEmailUniqChecker(repos UserRepository) *UserEmailUniqChecker {
	return &UserEmailUniqChecker{Repos: repos}
}

// IsUniqueEmail メールアドレスがユニークかどうかをチェックする。自身のメールアドレスは対象としないようにする
func (u *UserEmailUniqChecker) IsUniqueEmail(newUser *UserModel) (bool, error) {
	user, err := u.Repos.GetUserByEmail(newUser.Email)
	if err != nil {
		if ErrNotFound.Error() == err.Error() {
			return true, nil
		}
		return false, errors.WithStack(err)
	}

	return user.ID == newUser.ID, nil
}
