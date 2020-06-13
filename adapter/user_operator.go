package adapter

import (
	"clean-serverless-book-sample-v2/domain"
	"github.com/guregu/dynamo"
	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
)

// UserOperator ユーザーを操作する構造体
type UserOperator struct {
	Client                 *ResourceTableOperator
	Mapper                 *DynamoModelMapper
	UserEmailUniqGenerator *UserEmailUniqGenerator
}

func (u *UserOperator) getUserResourceByID(id uint64) (*UserResource, error) {
	var user UserResource
	_, err := u.Mapper.GetEntityByID(id, &UserResource{}, &user)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

// GetUserByEmail メールアドレスからユーザー情報を取得する
func (u *UserOperator) GetUserByEmail(email string) (*domain.UserModel, error) {
	table, err := u.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fb := nomof.NewBuilder()
	fb.Equal("Email", email)
	fb.BeginsWith("PK", u.Mapper.GetEntityNameFromStruct(UserResource{}))

	var usersDynamo []UserResource
	err = table.Scan().Filter(fb.JoinAnd(), fb.Arg...).All(&usersDynamo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(usersDynamo) == 0 {
		return nil, errors.WithStack(domain.ErrNotFound)
	}

	return &usersDynamo[0].UserModel, nil
}

// Execute IDからユーザー情報を取得する
func (u *UserOperator) GetUserByID(id uint64) (*domain.UserModel, error) {
	userResource, err := u.getUserResourceByID(id)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return nil, errors.WithStack(domain.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &userResource.UserModel, nil
}

// Execute ユーザー一覧を取得する
func (u *UserOperator) GetUsers() ([]*domain.UserModel, error) {
	table, err := u.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fb := nomof.NewBuilder()
	fb.BeginsWith("PK", u.Mapper.GetEntityNameFromStruct(UserResource{}))

	var userDynamo []UserResource
	err = table.Scan().Filter(fb.JoinAnd(), fb.Arg...).All(&userDynamo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var users = make([]*domain.UserModel, len(userDynamo))
	for i := range userDynamo {
		users[i] = &userDynamo[i].UserModel
	}

	return users, nil
}

// CreateUser ユーザーを新規作成する
func (u *UserOperator) CreateUser(userModel *domain.UserModel) (*domain.UserModel, error) {
	conn, err := u.Client.ConnectDB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResource := NewUserResource(userModel, u.Mapper)

	tx := conn.WriteTx()

	r, err := u.Mapper.BuildQueryCreate(userResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	uniq, err := u.UserEmailUniqGenerator.BuildQueryCreateByUser(userResource)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = tx.Put(r).Put(uniq).Run()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &userResource.UserModel, nil
}

// UpdateUser ユーザーを更新する
func (u *UserOperator) UpdateUser(newUserModel *domain.UserModel) error {
	conn, err := u.Client.ConnectDB()
	if err != nil {
		return errors.WithStack(err)
	}

	oldUserResource, err := u.getUserResourceByID(newUserModel.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	newUserResource := *oldUserResource
	newUserResource.Email = newUserModel.Email
	newUserResource.Name = newUserModel.Name

	tx := conn.WriteTx()

	r, err := u.Mapper.BuildQueryUpdate(&newUserResource)
	if err != nil {
		return errors.WithStack(err)
	}

	query := tx.Put(r)

	if oldUserResource.Email != newUserResource.Email {
		uniqDelete, err := u.UserEmailUniqGenerator.BuildQueryDeleteByUser(oldUserResource)
		if err != nil {
			return errors.WithStack(err)
		}

		uniqCreate, err := u.UserEmailUniqGenerator.BuildQueryCreateByUser(&newUserResource)
		if err != nil {
			return errors.WithStack(err)
		}

		query.
			Put(uniqCreate).
			Delete(uniqDelete)
	}

	err = query.Run()

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// DeleteUser ユーザー情報を削除する
func (u *UserOperator) DeleteUser(userModel *domain.UserModel) error {
	conn, err := u.Client.ConnectDB()
	if err != nil {
		return errors.WithStack(err)
	}

	tx := conn.WriteTx()

	userResource := NewUserResource(userModel, u.Mapper)

	r, err := u.Mapper.BuildQueryDelete(userResource)
	if err != nil {
		return errors.WithStack(err)
	}

	uniq, err := u.UserEmailUniqGenerator.BuildQueryDeleteByUser(userResource)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Delete(r).Delete(uniq).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
