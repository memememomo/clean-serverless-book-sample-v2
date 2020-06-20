package controller

import (
	"fmt"
	"gopkg.in/validator.v2"
)

// validateErrorMessages エラ〜メッセージテンプレート
var validateErrorMessages = map[error]string{
	validator.ErrUnsupported: "%sは不正な値です。",
	validator.ErrZeroValue:   "%sを入力してください。",
	validator.ErrLen:         "%sの文字列長が不正です。",
	validator.ErrMax:         "%sの文字数が上限を超えています。",
	ErrRequired:              "%sを入力してください。",
	ErrEmail:                 "%sの形式が不正です。",
	ErrUint:                  "%sは0以上の数値を入力してください。",
	ErrUniq:                  "すでに登録されている%sです。",
}

// displayNames 引数名の日本語表示
var displayNames = map[string]string{
	"user_id":      "ユーザーID",
	"user_name":    "ユーザー名",
	"micropost_id": "マイクロポストID",
	"email":        "メールアドレス",
	"content":      "本文",
	"name":         "名前",
}

// ConvertErrorsToMessage エラーメッセージに変換
func ConvertErrorsToMessage(errs map[string]error) map[string]string {
	messages := map[string]string{}

	for argName, err := range errs {
		disp := displayNames[argName]
		if disp == "" {

		}
		message := validateErrorMessages[err]
		if message == "" {
			messages[argName] = err.Error()
		} else {
			messages[argName] = fmt.Sprintf(message, disp)
		}
	}

	return messages
}
