package adapter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"time"
)

type DynamoResource interface {
	EntityName() string
	PK() string
	SetPK()
	SK() string
	SetSK()
	SetID(id uint64)
	ID() uint64
	SetVersion(v int)
	Version() int
	CreatedAt() time.Time
	SetCreatedAt(t time.Time)
	UpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

type DynamoModelMapper struct {
	Client    *ResourceTableOperator
	TableName string
	PKName    string
	SKName    string
}

func (d *DynamoModelMapper) GetEntityNameFromStruct(s interface{}) string {
	r := reflect.TypeOf(s)
	return r.Name()
}

func (d *DynamoModelMapper) BuildQueryCreate(resource DynamoResource) (*dynamo.Put, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := d.generateID(resource.EntityName())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resource.SetCreatedAt(time.Now())
	resource.SetUpdatedAt(time.Now())
	resource.SetID(id)
	resource.SetVersion(1)
	resource.SetPK()
	resource.SetSK()

	fb := nomof.NewBuilder()
	fb.AttributeNotExists(d.PKName)

	query := table.
		Put(resource).
		If(fb.JoinAnd(), fb.Arg...)

	return query, nil
}

func (d *DynamoModelMapper) BuildQueryUpdate(resource DynamoResource) (*dynamo.Put, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	oldVersion := resource.Version()

	resource.SetUpdatedAt(time.Now())
	resource.SetVersion(oldVersion + 1)

	fb := nomof.NewBuilder()
	fb.Equal("Version", oldVersion)

	query := table.
		Put(resource).
		If(fb.JoinAnd(), fb.Arg...)

	return query, nil
}

func (d *DynamoModelMapper) BuildQueryDelete(resource DynamoResource) (*dynamo.Delete, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query := table.
		Delete(d.PKName, resource.PK()).
		Range(d.SKName, resource.SK())

	return query, nil
}

func (d *DynamoModelMapper) CreateResource(resource DynamoResource) error {
	query, err := d.BuildQueryCreate(resource)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *DynamoModelMapper) UpdateResource(resource DynamoResource) error {
	query, err := d.BuildQueryUpdate(resource)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *DynamoModelMapper) DeleteResource(resource DynamoResource) error {
	query, err := d.BuildQueryDelete(resource)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *DynamoModelMapper) PutResource(resource DynamoResource) error {
	if d.isNewEntity(resource) {
		return d.CreateResource(resource)
	}
	return d.UpdateResource(resource)
}

func (d *DynamoModelMapper) GetPK(resource DynamoResource) string {
	return fmt.Sprintf("%s-%011d", resource.EntityName(), resource.ID())
}

func (d *DynamoModelMapper) GetSK(resource DynamoResource) string {
	return fmt.Sprintf("%011d", resource.ID())
}

func (d *DynamoModelMapper) GetEntityByID(id uint64, resource DynamoResource, ret interface{}) (interface{}, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resource.SetID(id)
	err = table.
		Get(d.PKName, resource.PK()).
		Range(d.SKName, dynamo.Equal, resource.SK()).
		One(ret)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ret, nil
}

func (d *DynamoModelMapper) isNewEntity(resource DynamoResource) bool {
	return resource.Version() == 0
}

func (d *DynamoModelMapper) generateID(tableName string) (uint64, error) {
	attr, err := d.atomicCount(fmt.Sprintf("AtomicCounter-%s", tableName), "AtomicCounter", "CurrentNumber", 1)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	nStr := aws.StringValue(attr.N)
	n, err := strconv.ParseUint(nStr, 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return n, nil
}

func (d *DynamoModelMapper) atomicCount(pk, sk, counterName string, value int) (*dynamodb.AttributeValue, error) {
	db, err := d.Client.ConnectDB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	output, err := db.Client().UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			counterName: {
				Action: aws.String("ADD"),
				Value: &dynamodb.AttributeValue{
					N: aws.String(fmt.Sprintf("%d", value)),
				},
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return output.Attributes[counterName], nil
}
