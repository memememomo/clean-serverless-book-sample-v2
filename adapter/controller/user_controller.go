package controller

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/interactor"
	"clean-serverless-book-sample-v2/registry"
	"clean-serverless-book-sample-v2/usecase"
	"clean-serverless-book-sample-v2/utils"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

// PostSettingValidator バリデーション設定
func PostSettingValidator() *Validator {
	return &Validator{
		Settings: []*ValidatorSetting{
			{ArgName: "user_name", ValidateTags: "required"},
			{ArgName: "email", ValidateTags: "required,email"},
		},
	}
}

// RequestPostUser PostUserのリクエスト
type RequestPostUser struct {
	Name  string `json:"user_name"`
	Email string `json:"email"`
}

// RequestPutUser PutUserのリクエスト
type RequestPutUser struct {
	Name  string `json:"user_name"`
	Email string `json:"email"`
}

// UserResponse レスポンス用のJSON形式を表した構造体
type UserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"user_name"`
	Email string `json:"email"`
}

// UserResponse Userリストレスポンス用のJSON形式を表した構造体
type UsersResponse struct {
	Users []*UserResponse `json:"users"`
}

// PostUsers 新規作成
func PostUsers(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション処理
	validator := PostSettingValidator()
	validErr := validator.ValidateBody(request.Body)
	if validErr != nil {
		return Response400(validErr)
	}

	// JSON形式から構造体に変換
	var req RequestPostUser
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	// 新規作成処理
	creator := registry.GetFactory().BuildCreateUser()
	res, err := creator.Execute(&usecase.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		if err.Error() == interactor.ErrUniqEmail.Error() {
			return Response400(map[string]error{
				"email": errors.New("すでに登録されているメールアドレスです。"),
			})
		}
		return Response500(err)
	}

	// 201レスポンス
	return Response201(res.GetUserID())
}

// PutUser 更新
func PutUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション処理
	validator := PostSettingValidator()
	validErr := validator.ValidateBody(request.Body)
	if validErr != nil {
		return Response400(validErr)
	}

	// JSON形式から構造体に変換
	var req RequestPutUser
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	// パスパラメータからユーザーIDを取得する
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// 更新処理
	updater := registry.GetFactory().BuildUpdateUser()
	_, err = updater.Execute(&usecase.UpdateUserRequest{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		if err.Error() == interactor.ErrUniqEmail.Error() {
			return Response400(map[string]error{
				"email": errors.New("すでに登録されているメールアドレスです。"),
			})
		}
		return Response500(err)
	}

	// 200レスポンス
	return Response200OK()
}

// GetUsers 一覧取得処理
func GetUsers(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// 一覧取得処理
	getter := registry.GetFactory().BuildGetUserList()
	res, err := getter.Execute(&usecase.GetUserListRequest{})
	if err != nil {
		return Response500(err)
	}

	// ドメインモデルからレスポンス用の構造体に詰め替える
	var resUsers = make([]*UserResponse, res.UserCount())
	for i, u := range res.Users {
		resUsers[i] = &UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	// レスポンス処理
	return Response200(&UsersResponse{
		Users: resUsers,
	})
}

// GetUser IDから取得
func GetUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからユーザーIDを取得する
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// ユーザー取得処理
	getter := registry.GetFactory().BuildGetUserByID()
	res, err := getter.Execute(&usecase.GetUserByIDRequest{UserID: userID})
	if err != nil {
		if err.Error() == domain.ErrNotFound.Error() {
			return Response404()
		}
		return Response500(err)
	}

	// ドメインモデルからレスポンス用構造体に詰め替えて、レスポンス
	return Response200(&UserResponse{
		ID:    res.User.ID,
		Name:  res.User.Name,
		Email: res.User.Email,
	})
}

// DeleteUser 削除処理
func DeleteUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからユーザーIDを取得する
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// 削除処理
	deleter := registry.GetFactory().BuildUserDeleter()
	_, err = deleter.Execute(&usecase.DeleteUserRequest{
		UserID: userID,
	})
	if err != nil {
		return Response500(err)
	}

	// レスポンス
	return Response200OK()
}
