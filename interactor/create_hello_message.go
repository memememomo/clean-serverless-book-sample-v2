package interactor

import (
	"clean-serverless-book-sample-v2/usecase"
	"fmt"
)

// CreateHelloMessage Helloメッセージ作成
type CreateHelloMessage struct {
}

// NewCreateHelloMessage CreateHelloMessageインスタンスを生成
func NewCreateHelloMessage() *CreateHelloMessage {
	return &CreateHelloMessage{}
}

// Execute 実行
func (c *CreateHelloMessage) Execute(req *usecase.CreateHelloMessageRequest) (*usecase.CreateHelloMessageResponse, error) {
	msg := fmt.Sprintf("Hello!%s", req.Name)
	return &usecase.CreateHelloMessageResponse{Message: msg}, nil
}
