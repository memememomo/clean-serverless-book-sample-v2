package adapter

import (
	"encoding/base64"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/pkg/errors"
)

// AWSKmsClient AWS SDKから KMS を利用して暗号化・復号化する
type AWSKmsClient struct {
	Client *kms.KMS
	KeyID  string
}

// NewAWSKmsClient AWSKmsClient インスタンスを生成
func NewAWSKmsClient() *AWSKmsClient {
	client := kms.New(
		session.Must(session.NewSession()),
		aws.NewConfig().WithRegion("ap-northeast-1"))
	return &AWSKmsClient{
		Client: client,
		KeyID:  os.Getenv("KMS_KEY_ID"),
	}
}

// Decrypt 復号化
func (a *AWSKmsClient) Decrypt(str string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decode")
	}
	res, err := a.Client.Decrypt(&kms.DecryptInput{
		CiphertextBlob: decodedBytes,
	})
	if err != nil {
		return "", errors.Wrap(err, "Failed to decrypt by KMS")
	}

	return string(res.Plaintext[:]), nil
}

// Encrypt 暗号化
func (a *AWSKmsClient) Encrypt(str string) (string, error) {
	res, err := a.Client.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(a.KeyID),
		Plaintext: []byte(str),
	})
	if err != nil {
		return "", errors.Wrap(err, "Failed to encrypt by KMS")
	}

	encodedString := base64.StdEncoding.EncodeToString(res.CiphertextBlob)

	return encodedString, nil
}
