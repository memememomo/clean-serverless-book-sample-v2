package controller

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/glog"
)

// Response201Body IDを含めた201レスポンス
type Response201Body struct {
	Message string `json:"message"`
	ID      uint64 `json:"id"`
}

// Response400Body バリデーションエラーメッセージを含めた400レスポンス
type Response400Body struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

// Response401Body 401レスポンス
type Response401Body struct {
	Message string `json:"message"`
}

// commonHeaders 各レスポンスに共通で含むヘッダー
func commonHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
}

// Response200 JSONを含めた200レスポンス
func Response200(body interface{}) events.APIGatewayProxyResponse {
	b, err := json.Marshal(body)
	if err != nil {
		return Response500(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
		Headers:    commonHeaders(),
	}
}

// Response200OK okメッセージを含めた200レスポンス
func Response200OK() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    commonHeaders(),
		Body:       `{"message":"OK"}`,
	}
}

// Response201 IDを含めた201レスポンス
func Response201(id uint64) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    commonHeaders(),
		Body:       fmt.Sprintf(`{"message":"OK","id":%d}`, id),
	}
}

// Response400 エラーメッセージを含めた400レスポンス
func Response400(errs map[string]error) events.APIGatewayProxyResponse {
	glog.Warningf("%+v", errs)
	res := &Response400Body{
		Message: "入力値を確認してください。",
		Errors:  ConvertErrorsToMessage(errs),
	}

	b, err := json.Marshal(res)
	if err != nil {
		return Response500(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers:    commonHeaders(),
		Body:       string(b),
	}
}

// Response404 404レスポンス
func Response404() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers:    commonHeaders(),
		Body:       `{"message":"結果が見つかりません。"}`,
	}
}

// Response500 500レスポンス
func Response500(err error) events.APIGatewayProxyResponse {
	glog.Errorf("%+v\n", err)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers:    commonHeaders(),
		Body:       `{"message":"サーバエラーが発生しました。"}`,
	}
}
