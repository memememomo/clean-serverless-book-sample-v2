package adapter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDynamoClient_Connect(t *testing.T) {
	client := NewClient(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
		Endpoint:    aws.String(os.Getenv("DYNAMO_LOCAL_ENDPOINT")),
	})

	svc, err := client.Connect()
	assert.NoError(t, err)

	ret, err := svc.ListTables().All()
	assert.NoError(t, err)
	assert.Len(t, ret, 0)
}
