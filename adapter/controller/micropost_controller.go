package controller

import (
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/registry"
	"clean-serverless-book-sample-v2/usecase"
	"clean-serverless-book-sample-v2/utils"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

// MicropostSettingsValidator バリデーション設定
func MicropostSettingsValidator() *Validator {
	return &Validator{
		Settings: []*ValidatorSetting{{ArgName: "content", ValidateTags: "required,max=140"}},
	}
}

// RequestMicropost HTTPリクエストで送られてくるJSON形式を表した構造体
type RequestMicropost struct {
	Content string `json:"content"`
}

// RequestPostMicropost PostMicropostのリクエスト
type RequestPostMicropost struct {
	RequestMicropost
}

// RequestPutMicropost PutMicropostのリクエスト
type RequestPutMicropost struct {
	RequestMicropost
}

// ResponseMicropost レスポンス用のJSON形式を表した構造体
type ResponseMicropost struct {
	ID      uint64 `json:"id"`
	UserID  uint64 `json:"user_id"`
	Content string `json:"content"`
}

// ResponseMicroposts Micropostリストレスポンス用のJSON形式を表した構造体
type ResponseMicroposts struct {
	Microposts []*ResponseMicropost `json:"microposts"`
}

// PostMicroposts 新規作成
func PostMicroposts(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション処理
	validator := MicropostSettingsValidator()
	validErr := validator.ValidateBody(request.Body)
	if validErr != nil {
		return Response400(validErr)
	}

	// パスパラメータからユーザーIDを取得する
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// JSON形式から構造体に変換
	var req RequestPostMicropost
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	// 新規作成処理
	creator := registry.GetFactory().BuildCreateMicropost()
	res, err := creator.Execute(&usecase.CreateMicropostRequest{
		Content: req.Content,
		UserID:  userID,
	})
	if err != nil {
		return Response500(err)
	}

	// 201レスポンス
	return Response201(res.MicropostID)
}

// PutMicropost 更新
func PutMicropost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション処理
	validator := MicropostSettingsValidator()
	validErr := validator.ValidateBody(request.Body)
	if validErr != nil {
		return Response400(validErr)
	}

	// パスパラメータからユーザーIDを取得する
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// パスパラメータからマイクロポストIDを取得する
	micropostID, err := utils.ParseUint(request.PathParameters["micropost_id"])
	if err != nil {
		return Response500(err)
	}

	// JSON形式から構造体に変換
	var req RequestPutMicropost
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	// 更新処理
	updater := registry.GetFactory().BuildUpdateMicropost()
	_, err = updater.Execute(&usecase.UpdateMicropostRequest{
		Content:     req.Content,
		UserID:      userID,
		MicropostID: micropostID,
	})
	if err != nil {
		return Response500(err)
	}

	// 200レスポンス
	return Response200OK()
}

// Execute 一覧取得
func GetMicroposts(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからユーザーIDを取得
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// マイクロポスト取得処理
	getter := registry.GetFactory().BuildGetMicropostList()
	res, err := getter.Execute(&usecase.GetMicropostListRequest{
		UserID: userID,
	})
	if err != nil {
		return Response500(err)
	}

	// ドメインモデルからレスポンス用の構造体に詰め替える
	var resMicroposts = make([]*ResponseMicropost, len(res.Microposts))
	for i, m := range res.Microposts {
		resMicroposts[i] = &ResponseMicropost{
			ID:      m.ID,
			UserID:  m.UserID,
			Content: m.Content,
		}
	}

	// レスポンス処理
	return Response200(&ResponseMicroposts{
		Microposts: resMicroposts,
	})
}

// GetMicropost IDから取得
func GetMicropost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからユーザーIDを取得する
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// パスパラメータからマイクロポストIDを取得する
	micropostID, err := utils.ParseUint(request.PathParameters["micropost_id"])
	if err != nil {
		return Response500(err)
	}

	// マイクロポスト取得処理
	getter := registry.GetFactory().BuildGetMicropostByID()
	res, err := getter.Execute(&usecase.GetMicropostByIDRequest{
		MicropostID: micropostID,
		UserID:      userID,
	})
	if err != nil {
		if err.Error() == domain.ErrNotFound.Error() {
			return Response404()
		}
		return Response500(err)
	}

	// ドメインモデルからレスポンス用構造体に詰め替えて、レスポンス
	return Response200(&ResponseMicropost{
		ID:      res.Micropost.ID,
		Content: res.Micropost.Content,
		UserID:  res.Micropost.UserID,
	})
}

// DeleteMicropost 削除処理
func DeleteMicropost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからユーザーIDを取得する
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	// パスパラメータからマイクロポストIDを取得する
	micropostID, err := utils.ParseUint(request.PathParameters["micropost_id"])
	if err != nil {
		return Response500(err)
	}

	// 削除処理
	deleter := registry.GetFactory().BuildDeleteMicropost()
	_, err = deleter.Execute(&usecase.DeleteMicropostRequest{
		MicropostID: micropostID,
		UserID:      userID,
	})
	if err != nil {
		return Response500(err)
	}

	// レスポンス
	return Response200OK()
}
