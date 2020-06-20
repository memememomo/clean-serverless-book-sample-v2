package controller

import (
	"clean-serverless-book-sample/mocks"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestPostHello_200 リクエストが成功する場合
func TestPostHello_200(t *testing.T) {
	// DynamoDB Local の設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// テスト用のリクエストBoyd
	body := map[string]interface{}{
		"name": "Taro",
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	// テスト用のリクエスト
	req := events.APIGatewayProxyRequest{
		Body: string(bodyStr),
	}

	// 実行呼び出し
	res := PostHello(req)

	// レスポンスをmap形式に変換
	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	// ステータスコードを確認
	assert.Equal(t, 200, res.StatusCode)

	// メッセージを確認
	assert.Equal(t, "Hello!Taro", resBody["message"])
}

// TestPostHello_400 バリデーションエラーが発生する場合
func TestPostHello_400(t *testing.T) {
	// DynamoDB Localの設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// テスト用のリクエストBody
	body := map[string]interface{}{
		"name": "",
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	// テスト用のリクエスト
	req := events.APIGatewayProxyRequest{
		Body: string(bodyStr),
	}

	// 実行呼び出し
	res := PostHello(req)

	// レスポンスをmap形式に変換
	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	// ステータスコードを確認
	assert.Equal(t, 400, res.StatusCode)

	// エラーメッセージを確認
	errs := resBody["errors"].(map[string]interface{})
	assert.Equal(t, "名前を入力してください。", errs["name"])
}
