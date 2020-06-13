package controller

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/mocks"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestPostUsers_201 新規作成 成功時
func TestPostUsers_201(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// リクエストパラメータ設定
	body := map[string]interface{}{
		"user_name": "テスト名前",
		"email":     "test@example.com",
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	// 新規作成処理
	res := PostUsers(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
	})

	// レスポンスコードをチェック
	assert.Equal(t, 201, res.StatusCode)

	// JSONからmap型に変換
	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	// IDをチェック
	id := uint64(resBody["id"].(float64))
	assert.Equal(t, uint64(1), id)

	// DynamoDBに保存されたデータをチェック
	user, err := tables.UserOperator.GetUserByID(id)
	assert.NoError(t, err)
	assert.Equal(t, body["user_name"].(string), user.Name)
	assert.Equal(t, body["email"].(string), user.Email)
}

// TestPostUsers_400 新規登録 バリデーションエラー時
func TestPostUsers_400(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 重複エラーテスト用のモックデータを作成
	userMock, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    1,
		Name:  "Name_1",
		Email: "test1@example.com",
	})
	assert.NoError(t, err)

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		// 未入力の場合
		{
			Request: map[string]interface{}{
				"user_name": "",
				"email":     "",
			},
			Expected: map[string]interface{}{
				"user_name": "ユーザー名を入力してください。",
				"email":     "メールアドレスを入力してください。",
			},
		},
		// メールアドレスの形式が不正の場合
		{
			Request: map[string]interface{}{
				"user_name": "hoge",
				"email":     "test@",
			},
			Expected: map[string]interface{}{
				"email": "メールアドレスの形式が不正です。",
			},
		},
		// メールアドレス重複の場合
		{
			Request: map[string]interface{}{
				"user_name": "dup",
				"email":     userMock.Email,
			},
			Expected: map[string]interface{}{
				"email": "すでに登録されているメールアドレスです。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PostUsers(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
		})

		var resBody map[string]interface{}
		err = json.Unmarshal([]byte(res.Body), &resBody)
		assert.NoError(t, err)

		errors := resBody["errors"].(map[string]interface{})

		assert.Equal(t, 400, res.StatusCode, msg)
		assert.Equal(t, c.Expected, errors, msg)
	}
}

// TestPutUser_200 更新 正常時
func TestPutUser_200(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 更新用モックデータを作成
	userMock, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    1,
		Name:  "Name_1",
		Email: "test1@example.com",
	})
	assert.NoError(t, err)

	// 更新リクエストパラメータ
	body := map[string]interface{}{
		"user_name": "テスト名前更新",
		"email":     "test_update@example.com",
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	// 更新処理
	res := PutUser(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	// レスポンスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// JSONからmap型に変換
	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	// DynamoDBのデータが更新されているかをチェック
	user, err := tables.UserOperator.GetUserByID(userMock.ID)
	assert.NoError(t, err)
	assert.Equal(t, body["user_name"].(string), user.Name)
	assert.Equal(t, body["email"].(string), user.Email)
}

// TestPutUser_200_dup 更新 メールアドレスを変更しない場合(重複エラーにならないこと)
func TestPutUser_200_dup(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// モックデータを作成
	userMock, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    1,
		Name:  "Name_1",
		Email: "test1@example.com",
	})
	assert.NoError(t, err)

	// 更新パラメータ
	body := map[string]interface{}{
		"user_name": "テスト名前更新",
		"email":     userMock.Email,
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	// 更新処理
	res := PutUser(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	// ステータスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// JSONからmap型に変換
	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	// DynamoDBのデータをチェック
	user, err := tables.UserOperator.GetUserByID(userMock.ID)
	assert.NoError(t, err)
	assert.Equal(t, body["user_name"].(string), user.Name)
	assert.Equal(t, body["email"].(string), user.Email)
}

// TestPutUser_400 更新 バリデーションエラー時
func TestPutUser_400(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 更新用モックデータを作成
	userMock, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    1,
		Name:  "Name_1",
		Email: "test1@example.com",
	})
	assert.NoError(t, err)

	// 重複エラー用モックデータを作成
	dupUserMock, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    2,
		Name:  "Name_2",
		Email: "test1@cample.com",
	})
	assert.NoError(t, err)

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		// 未入力の場合
		{
			Request: map[string]interface{}{
				"user_name": "",
				"email":     "",
			},
			Expected: map[string]interface{}{
				"user_name": "ユーザー名を入力してください。",
				"email":     "メールアドレスを入力してください。",
			},
		},
		// メールアドレスの形式が不正な場合
		{
			Request: map[string]interface{}{
				"user_name": "hoge",
				"email":     "test@",
			},
			Expected: map[string]interface{}{
				"email": "メールアドレスの形式が不正です。",
			},
		},
		// 重複エラーの場合
		{
			Request: map[string]interface{}{
				"user_name": "hoge",
				"email":     dupUserMock.Email,
			},
			Expected: map[string]interface{}{
				"email": "すでに登録されているメールアドレスです。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PutUser(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
			PathParameters: map[string]string{
				"user_id": fmt.Sprintf("%d", userMock.ID),
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

// TestGetUser 取得 正常時
func TestGetUser(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 取得用モックデータを作成
	userMock, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    1,
		Name:  "Name_1",
		Email: "test1@example.com",
	})
	assert.NoError(t, err)

	// 取得処理
	res := GetUser(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	// ステータスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// 取得したデータをチェック
	var body map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)
	assert.Equal(t, float64(userMock.ID), body["id"])
	assert.Equal(t, userMock.Name, body["user_name"])
	assert.Equal(t, userMock.Email, body["email"])
}

// TestGetUsers 一覧取得
func TestGetUsers(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 取得用モックデータを作成
	userMock1, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    1,
		Name:  "Name_1",
		Email: "test1@example.com",
	})
	assert.NoError(t, err)

	userMock2, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    2,
		Name:  "Name_2",
		Email: "test2@example.com",
	})

	// 一覧取得処理
	res := GetUsers(events.APIGatewayProxyRequest{})

	// ステータスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// JSONからmap型に変換
	var body map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)

	// 取得したデータをチェック
	actualUsers := body["users"].([]interface{})

	expected1 := userMock2
	actual1 := actualUsers[0].(map[string]interface{})
	assert.Equal(t, float64(expected1.ID), actual1["id"])
	assert.Equal(t, expected1.Name, actual1["user_name"])
	assert.Equal(t, expected1.Email, actual1["email"])

	expected2 := userMock1
	actual2 := actualUsers[1].(map[string]interface{})
	assert.Equal(t, float64(expected2.ID), actual2["id"])
	assert.Equal(t, expected2.Name, actual2["user_name"])
	assert.Equal(t, expected2.Email, actual2["email"])
}

// TestDeleteUser 削除
func TestDeleteUser(t *testing.T) {
	// テスト用DynamoDBを設定
	tables := mocks.SetupDB(t)
	defer tables.Cleanup()

	// 削除用モックデータを作成
	userMock, err := tables.UserOperator.CreateUser(&domain.UserModel{
		ID:    1,
		Name:  "Name_1",
		Email: "test1@example.com",
	})
	assert.NoError(t, err)

	// 削除処理
	res := DeleteUser(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	// ステータスコードをチェック
	assert.Equal(t, 200, res.StatusCode)

	// DynamoDBからデータが削除されているかをチェック
	users, err := tables.UserOperator.GetUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}
