package registry

import (
	"clean-serverless-book-sample-v2/adapter"
	"clean-serverless-book-sample-v2/domain"
	"clean-serverless-book-sample-v2/interactor"
	"clean-serverless-book-sample-v2/usecase"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// FactorySingleton Factoryのインスタンスを使い回すための変数
var FactorySingleton *Factory

// Factory 様々なインスタンスを生成する構造体
type Factory struct {
	Envs  *Envs
	cache map[string]interface{}
}

// ClearFactory 使いまわしているインスタンスを削除する
func ClearFactory() {
	FactorySingleton = nil
}

// GetFactory Factoryのインスタンスを取得する
func GetFactory() *Factory {
	if FactorySingleton == nil {
		FactorySingleton = &Factory{
			Envs: NewEnvs(),
		}
	}
	return FactorySingleton
}

// container cacheにインスタンスがある場合はそれを返し、なければ新規作成する
func (f *Factory) container(key string, builder func() interface{}) interface{} {
	if f.cache == nil {
		f.cache = map[string]interface{}{}
	}
	if f.cache[key] == nil {
		f.cache[key] = builder()
	}
	return f.cache[key]
}

// BuildDynamoClient DynamoDBに接続するためのインスタンスを生成
func (f *Factory) BuildDynamoClient() *adapter.DynamoClient {
	return f.container("DynamoClient", func() interface{} {
		config := &aws.Config{
			Region: aws.String("ap-northeast-1"),
		}

		if f.Envs.DynamoLocalEndpoint() != "" {
			config.Credentials = credentials.NewStaticCredentials("dummy", "dummy", "dummy")
			config.Endpoint = aws.String(f.Envs.DynamoLocalEndpoint())
		}
		return adapter.NewClient(config)
	}).(*adapter.DynamoClient)
}

// BuildResourceTableOperator DynamoDBのテーブルに接続するためのインスタンスを生成
func (f *Factory) BuildResourceTableOperator() *adapter.ResourceTableOperator {
	return f.container("ResourceTableOperator", func() interface{} {
		return adapter.NewResourceTableOperator(
			f.BuildDynamoClient(),
			f.Envs.DynamoTableName())
	}).(*adapter.ResourceTableOperator)
}

// BuildDynamoModelMapper ModelからDynamoDBに保存する形式に変換するためのインスタンスを生成
func (f *Factory) BuildDynamoModelMapper() *adapter.DynamoModelMapper {
	return f.container("DynamoModelMapper", func() interface{} {
		return &adapter.DynamoModelMapper{
			Client:    f.BuildResourceTableOperator(),
			TableName: f.Envs.DynamoTableName(),
			PKName:    f.Envs.DynamoPKName(),
			SKName:    f.Envs.DynamoSKName(),
		}
	}).(*adapter.DynamoModelMapper)
}

// BuildUserEmailUniqGenerator ユーザーのメールアドレス重複チェック用のレコード生成機のインスタンスを生成
func (f *Factory) BuildUserEmailUniqGenerator() *adapter.UserEmailUniqGenerator {
	return f.container("UserEmailUniqGenerator", func() interface{} {
		return adapter.NewUserEmailUniqGenerator(
			f.BuildDynamoModelMapper(),
			f.BuildResourceTableOperator(),
			f.Envs.DynamoPKName(),
			f.Envs.DynamoSKName())
	}).(*adapter.UserEmailUniqGenerator)
}

// BuildUserOperator ユーザー情報関連の操作を行うインスタンスを生成
func (f *Factory) BuildUserOperator() domain.UserRepository {
	return f.container("UserOperator", func() interface{} {
		return &adapter.UserOperator{
			Client:                 f.BuildResourceTableOperator(),
			Mapper:                 f.BuildDynamoModelMapper(),
			UserEmailUniqGenerator: f.BuildUserEmailUniqGenerator(),
		}
	}).(domain.UserRepository)
}

// BuildUserEmailUniqChecker ユーザーのメールアドレス重複チェックインスタンスを生成
func (f *Factory) BuildUserEmailUniqChecker() *domain.UserEmailUniqChecker {
	return f.container("UserEmailUniqChecker", func() interface{} {
		return domain.NewUserEmailUniqChecker(f.BuildUserOperator())
	}).(*domain.UserEmailUniqChecker)
}

// BuildMicropostOperator マイクロポスト情報関連の操作を行うインスタンスを生成
func (f *Factory) BuildMicropostOperator() *adapter.MicropostOperator {
	return f.container("MicropostOperator", func() interface{} {
		return &adapter.MicropostOperator{
			Client: f.BuildResourceTableOperator(),
			Mapper: f.BuildDynamoModelMapper(),
		}
	}).(*adapter.MicropostOperator)
}

// BuildCreateUser ユーザー作成UseCaseインスタンスを生成
func (f *Factory) BuildCreateUser() usecase.ICreateUser {
	return f.container("CreateUser", func() interface{} {
		return interactor.NewCreateUser(
			f.BuildUserOperator(),
			f.BuildUserEmailUniqChecker())
	}).(usecase.ICreateUser)
}

// BuildUpdateUser ユーザー更新UseCaseインスタンスを生成
func (f *Factory) BuildUpdateUser() usecase.IUpdateUser {
	return f.container("UpdateUser", func() interface{} {
		return interactor.NewUpdateUser(
			f.BuildUserOperator(),
			f.BuildUserEmailUniqChecker())
	}).(usecase.IUpdateUser)
}

// BuildGetUserList ユーザー取得UseCaseインスタンスを生成
func (f *Factory) BuildGetUserList() usecase.IGetUserList {
	return f.container("GetUserList", func() interface{} {
		return interactor.NewGetUserList(f.BuildUserOperator())
	}).(usecase.IGetUserList)
}

// BuildGetUserByID ユーザー取得UseCaseインスタンスを生成
func (f *Factory) BuildGetUserByID() usecase.IGetUserByID {
	return f.container("Execute", func() interface{} {
		return interactor.NewGetUserByID(f.BuildUserOperator())
	}).(usecase.IGetUserByID)
}

// BuildUserDeleter ユーザー削除Usecaseインスタンスを生成
func (f *Factory) BuildUserDeleter() usecase.IDeleteUser {
	return f.container("UserDeleter", func() interface{} {
		return interactor.NewUserDeleter(
			f.BuildUserOperator(),
			f.BuildGetUserByID())
	}).(usecase.IDeleteUser)
}

// BuildCreateMicropost マイクロポスト作成UseCaseインスタンスを生成
func (f *Factory) BuildCreateMicropost() usecase.ICreateMicropost {
	return f.container("CreateMicropost", func() interface{} {
		return interactor.NewCreateMicropost(
			f.BuildMicropostOperator())
	}).(usecase.ICreateMicropost)
}

// BuildGetMicropostList マイクロポスト取得UseCaseインスタンスを生成
func (f *Factory) BuildGetMicropostList() usecase.IGetMicropostList {
	return f.container("GetMicropostList", func() interface{} {
		return interactor.NewGetMicropostList(
			f.BuildMicropostOperator())
	}).(usecase.IGetMicropostList)
}

// BuildGetMicropostList マイクロポスト取得UseCaseインスタンスを生成
func (f *Factory) BuildGetMicropostByID() usecase.IGetMicropostByID {
	return f.container("GetMicropostByID", func() interface{} {
		return interactor.NewGetMicropostByID(
			f.BuildMicropostOperator())
	}).(usecase.IGetMicropostByID)
}

// BuildUpdateMicropost マイクロポスト更新UseCaseインスタンスを生成
func (f *Factory) BuildUpdateMicropost() usecase.IUpdateMicropost {
	return f.container("UpdateMicropost", func() interface{} {
		return interactor.NewUpdateMicropost(
			f.BuildMicropostOperator())
	}).(usecase.IUpdateMicropost)
}

// BuildDeleteMicropost マイクロポスト削除UseCaseインスタンスを生成
func (f *Factory) BuildDeleteMicropost() usecase.IDeleteMicropost {
	return f.container("DeleteMicropost", func() interface{} {
		return interactor.NewDeleteMicropost(
			f.BuildGetMicropostByID(),
			f.BuildMicropostOperator())
	}).(usecase.IDeleteMicropost)
}

func (f *Factory) BuildCreateHelloMessage() usecase.ICreateHelloMessage {
	return f.container("CreateHelloMessage", func() interface{} {
		return interactor.NewCreateHelloMessage()
	}).(usecase.ICreateHelloMessage)
}
