package registry

import (
	"clean-serverless-book-sample-v2/adapter"
	"github.com/golang/glog"
	"os"
)

// Envs 環境変数を扱う。暗号化やキャッシュなどもできるようになっている
type Envs struct {
	KMSClient *adapter.AWSKmsClient
	Cache     map[string]string
}

var envs *Envs

// NewEnvs Envs インスタンスを生成
func NewEnvs() *Envs {
	return &Envs{
		KMSClient: adapter.NewAWSKmsClient(),
		Cache:     make(map[string]string),
	}
}

// Env シングルトンを取得する
func Env() *Envs {
	if envs == nil {
		envs = NewEnvs()
	}
	return envs
}

func (c *Envs) decrypt(key string) string {
	if os.Getenv("DISABLE_ENV_DECRYPT") != "" {
		return c.env(key)
	}

	v := c.Cache[key]
	if v != "" {
		return v
	}

	str := os.Getenv(key)
	if str == "" {
		return ""
	}

	v, err := c.KMSClient.Decrypt(str)
	if err != nil {
		glog.Warning(err.Error())
		return ""
	}

	c.Cache[key] = v

	return c.Cache[key]
}

func (c *Envs) env(key string) string {
	return os.Getenv(key)
}

func (c *Envs) DynamoLocalEndpoint() string {
	return c.env("DYNAMO_LOCAL_ENDPOINT")
}

func (c *Envs) DynamoTableName() string {
	return c.env("DYNAMO_TABLE_NAME")
}

func (c *Envs) DynamoPKName() string {
	return c.env("DYNAMO_PK_NAME")
}

func (c *Envs) DynamoSKName() string {
	return c.env("DYNAMO_SK_NAME")
}
