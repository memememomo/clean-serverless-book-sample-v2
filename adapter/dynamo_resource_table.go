package adapter

import (
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

type Resource interface {
	SetPK()
	GetPK() string
}

type ResourceTableOperator struct {
	TableOperator
}

type ResourceSchema struct {
	PK string `dynamo:"PK,hash"`
	SK string `dynamo:"SK,range"`
}

func NewResourceTableOperator(client *DynamoClient, tableName string) *ResourceTableOperator {
	return &ResourceTableOperator{
		TableOperator: *NewTableOperator(client, tableName),
	}
}

func (r *ResourceTableOperator) getFromDynamo(query Resource, ret interface{}) error {
	table, err := r.ConnectTable()
	if err != nil {
		return errors.WithStack(err)
	}

	err = table.Get("PK", query.GetPK()).One(ret)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ResourceTableOperator) putToDynamo(query Resource) error {
	query.SetPK()

	table, err := r.ConnectTable()
	if err != nil {
		return errors.WithStack(err)
	}

	err = table.Put(query).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ResourceTableOperator) CreateTableForTest() error {
	return r.TableOperator.CreateTableForTest(&ResourceSchema{})
}

func (r *ResourceTableOperator) Dump() error {
	table, err := r.ConnectTable()
	if err != nil {
		return errors.WithStack(err)
	}

	var data []map[string]interface{}
	err = table.Scan().All(&data)
	if err != nil {
		return errors.WithStack(err)
	}

	pp.Print(data)

	return nil
}
