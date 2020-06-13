package mocks

import (
	"clean-serverless-book-sample-v2/adapter"
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/registry"
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

type DynamoTableOperator struct {
	Operator          *adapter.ResourceTableOperator
	UserOperator      domain.UserRepository
	MicropostOperator domain.MicropostRepository
}

func SetupDB(t *testing.T) *DynamoTableOperator {
	t.Helper()

	os.Setenv("DYNAMO_TABLE_NAME", generateRandomTableName(t))

	registry.ClearFactory()
	f := registry.GetFactory()
	operator := &DynamoTableOperator{}
	operator.Operator = f.BuildResourceTableOperator()
	operator.UserOperator = f.BuildUserOperator()
	operator.MicropostOperator = f.BuildMicropostOperator()

	operator.Operator.CreateTableForTest()

	return operator
}

func (d *DynamoTableOperator) Cleanup() {
	d.Operator.DropTable()
}

func generateRandomTableName(t *testing.T) string {
	length := 60
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		t.Fatal(err)
	}
	l := length / 2
	if length%2 == 1 {
		l++
	}
	return fmt.Sprintf("%x", buf[0:l])
}
