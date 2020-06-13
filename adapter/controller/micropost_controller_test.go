package controller

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/mocks"
	"clean-serverless-book-sample-v2/registry"
	"clean-serverless-book-sample-v2/usecase"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestPostMicroposts_201 新規作成処理 正常時
func TestPostMicroposts_201(t *testing.T) {
	// テスト用DynamoDBの設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// リクエスト用パラメータ
	body := map[string]interface{}{
		"content": strings.Repeat("a", 140),
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	userID := uint64(1)

	// 新規作成処理
	res := PostMicroposts(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	})

	// レスポンスコードチェック
	assert.Equal(t, 201, res.StatusCode)

	// JSONからmap型に変換
	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	// 新規作成されたIDの値をチェック
	id := uint64(resBody["id"].(float64))
	assert.Equal(t, uint64(1), id)

	// DynamoDBに保存されているかチェック
	getter := registry.GetFactory().BuildGetMicropostByID()
	micropostRes, err := getter.Execute(&usecase.GetMicropostByIDRequest{
		UserID:      userID,
		MicropostID: id,
	})
	assert.NoError(t, err)
	assert.Equal(t, body["content"].(string), micropostRes.Micropost.Content)
	assert.Equal(t, userID, micropostRes.Micropost.UserID)
}

// TestPostMicroposts_400 新規作成処理 バリデーションエラー時
func TestPostMicroposts_400(t *testing.T) {
	// テスト用DynamoDB設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		// 未入力の場合
		{
			Request: map[string]interface{}{
				"content": "",
			},
			Expected: map[string]interface{}{
				"content": "本文を入力してください。",
			},
		},
		// 本文の文字数が上限を超えている場合
		{
			Request: map[string]interface{}{
				"content": strings.Repeat("a", 141),
			},
			Expected: map[string]interface{}{
				"content": "本文の文字数が上限を超えています。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PostMicroposts(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
			PathParameters: map[string]string{
				"user_id": "1",
			},
		})

		var resBody map[string]interface{}
		err = json.Unmarshal([]byte(res.Body), &resBody)
		assert.NoError(t, err)

		errors := resBody["errors"].(map[string]interface{})

		assert.Equal(t, 400, res.StatusCode, msg)
		assert.Equal(t, c.Expected, errors)
	}
}

// TestPutMicropost_200 更新処理 正常
func TestPutMicropost_200(t *testing.T) {
	// テスト用のDynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 更新用モックデータを作成
	micropostMock, err := tables.MicropostOperator.CreateMicropost(&domain.MicropostModel{
		Content: "Content_1",
		UserID:  1,
	})
	assert.NoError(t, err)

	// 更新用リクエスト
	body := map[string]interface{}{
		"content": strings.Repeat("a", 140),
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	// 更新処理
	res := PutMicropost(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
			"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
		},
	})

	// レスポンスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// JSONからmap型に変換
	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	// DynamoDBに更新データが反映されているかチェック
	micropost, err := tables.MicropostOperator.GetMicropostByID(micropostMock.ID)
	assert.NoError(t, err)
	assert.Equal(t, body["content"].(string), micropost.Content)
}

// TestPutMicropost_400 更新処理 バリデーションエラー時
func TestPutMicropost_400(t *testing.T) {
	// テスト用のDynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 更新用モックデータを作成
	micropostMock, err := tables.MicropostOperator.CreateMicropost(&domain.MicropostModel{
		Content: "Content_1",
		UserID:  1,
	})
	assert.NoError(t, err)

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		// 未入力時
		{
			Request: map[string]interface{}{
				"content": "",
			},
			Expected: map[string]interface{}{
				"content": "本文を入力してください。",
			},
		},
		// 文字数が上限を超えている場合
		{
			Request: map[string]interface{}{
				"content": strings.Repeat("a", 141),
			},
			Expected: map[string]interface{}{
				"content": "本文の文字数が上限を超えています。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PutMicropost(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
			PathParameters: map[string]string{
				"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
				"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
			},
		})

		var resBody map[string]interface{}
		err = json.Unmarshal([]byte(res.Body), &resBody)
		assert.NoError(t, err)

		errors := resBody["errors"].(map[string]interface{})

		assert.Equal(t, 400, res.StatusCode, msg)
		assert.Equal(t, c.Expected, errors)
	}
}

// TestGetMicropost 取得処理
func TestGetMicropost(t *testing.T) {
	// テスト用のDynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 取得用のモックデータを作成
	micropostMock, err := tables.MicropostOperator.CreateMicropost(&domain.MicropostModel{
		Content: "Content_1",
		UserID:  1,
	})
	assert.NoError(t, err)

	// 取得処理
	res := GetMicropost(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
			"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
		},
	})

	// 取得したデータをチェック
	var body map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)
	assert.Equal(t, float64(micropostMock.ID), body["id"])
	assert.Equal(t, micropostMock.Content, body["content"])
	assert.Equal(t, float64(micropostMock.UserID), body["user_id"])
}

// TestGetMicroposts 一覧取得処理
func TestGetMicroposts(t *testing.T) {
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 取得用のモックデータを作成
	micropostMock1, err := tables.MicropostOperator.CreateMicropost(&domain.MicropostModel{
		Content: "Content_1",
		UserID:  1,
	})
	assert.NoError(t, err)

	micropostMock2, err := tables.MicropostOperator.CreateMicropost(&domain.MicropostModel{
		Content: "Content_2",
		UserID:  1,
	})
	assert.NoError(t, err)

	// このデータはUserIDが異なるので取得されない想定
	_, err = tables.MicropostOperator.CreateMicropost(&domain.MicropostModel{
		Content: "Content_3",
		UserID:  2,
	})
	assert.NoError(t, err)

	// 一覧取得処理
	res := GetMicroposts(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id": "1",
		},
	})

	// レスポンスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// JSONからmap型に変換
	var body map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)

	// 取得したデータをチェック
	actualMicroposts := body["microposts"].([]interface{})

	expected1 := micropostMock1
	actual1 := actualMicroposts[0].(map[string]interface{})
	assert.Equal(t, float64(expected1.ID), actual1["id"])
	assert.Equal(t, expected1.Content, actual1["content"])
	assert.Equal(t, float64(expected1.UserID), actual1["user_id"])

	expected2 := micropostMock2
	actual2 := actualMicroposts[1].(map[string]interface{})
	assert.Equal(t, float64(expected2.ID), actual2["id"])
	assert.Equal(t, expected2.Content, actual2["content"])
	assert.Equal(t, float64(expected2.UserID), actual2["user_id"])
}

// TestDeleteMicropost 削除処理
func TestDeleteMicropost(t *testing.T) {
	// テスト用のDynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 削除用モックデータを作成
	micropostMock, err := tables.MicropostOperator.CreateMicropost(&domain.MicropostModel{
		Content: "Content_1",
		UserID:  1,
	})
	assert.NoError(t, err)

	// 削除処理
	res := DeleteMicropost(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
			"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
		},
	})

	// ステータスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// DynamoDBからデータが削除されているかチェック
	microposts, err := tables.MicropostOperator.GetMicropostsByUserID(micropostMock.UserID)
	assert.NoError(t, err)
	assert.Len(t, microposts, 0)
}
