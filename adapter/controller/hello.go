package controller

import (
	"clean-serverless-book-sample-v2/registry"
	"clean-serverless-book-sample-v2/usecase"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

// PostHelloRequest HTTPリクエストのJSON形式を表した構造体
type PostHelloRequest struct {
	Name string `json:"name"`
}

// HelloMessageResponse HTTPレスポンスのJSON形式を表した構造体
type HelloMessageResponse struct {
	Message string `json:"message"`
}

// バリデーションの設定
var ValidateHelloMessageSettings = []*ValidatorSetting{
	{ArgName: "name", ValidateTags: "required"},
}

// PostHello コントローラの実装
func PostHello(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション
	validErr := ValidateBody(request.Body, ValidateHelloMessageSettings)
	if validErr != nil {
		return Response400(validErr)
	}

	// JSONから構造体に変換
	var req PostHelloRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	// UseCaseを実行
	h := registry.GetFactory().BuildCreateHelloMessage()
	res, err := h.Execute(&usecase.CreateHelloMessageRequest{
		Name: req.Name,
	})
	if err != nil {
		return Response500(err)
	}

	// HTTPレスポンスを返す
	return Response200(&HelloMessageResponse{
		Message: res.Message,
	})
}
