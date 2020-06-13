package adapter

import (
	"github.com/guregu/dynamo"
	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
)

// UserEmailUniq メールアドレス重複チェック用のレコードを表した構造体
type UserEmailUniq struct {
	Email      string `dynamo:"PK"`
	EntityName string `dynamo:"SK"`
	Exists     bool   `dynamo:"Exists"`
	UserID     uint64 `dynamo:"UserID"`
}

type UserEmailUniqGenerator struct {
	Mapper *DynamoModelMapper
	Client *ResourceTableOperator
	PKName string
	SKName string
}

func NewUserEmailUniqGenerator(mapper *DynamoModelMapper, client *ResourceTableOperator, pkName, skName string) *UserEmailUniqGenerator {
	return &UserEmailUniqGenerator{
		Mapper: mapper,
		Client: client,
		PKName: pkName,
		SKName: skName,
	}
}

func (u *UserEmailUniqGenerator) NewUserEmailUniqByUser(user *UserResource) *UserEmailUniq {
	return &UserEmailUniq{
		Email:      user.Email,
		EntityName: u.Mapper.GetEntityNameFromStruct(*user),
		Exists:     true,
		UserID:     user.ID(),
	}
}

func (u *UserEmailUniqGenerator) BuildQueryCreateByUser(user *UserResource) (*dynamo.Put, error) {
	table, err := u.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	uniq := u.NewUserEmailUniqByUser(user)

	fb := nomof.NewBuilder()
	fb.AttributeNotExists("Exists")
	fb.Equal("UserID", user.ID())

	query := table.
		Put(uniq).
		If(fb.JoinOr(), fb.Arg...)

	return query, nil
}

func (u *UserEmailUniqGenerator) BuildQueryDeleteByUser(user *UserResource) (*dynamo.Delete, error) {
	table, err := u.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	uniq := u.NewUserEmailUniqByUser(user)

	query := table.
		Delete(u.PKName, uniq.Email).
		Range(u.SKName, uniq.EntityName)

	return query, nil
}
