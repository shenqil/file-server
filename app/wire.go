//go:build wireinject
// +build wireinject

package app

import (
	"fileServer/app/api"
	"fileServer/app/model/minio/bucket"
	"fileServer/app/router"
	"fileServer/app/service"

	"github.com/google/wire"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitAuth,
		InitMinio,
		bucket.BucketSet,
		InjectorSet,
		InitGinEngine,
		router.RouterSet,
		service.ServiceSet,
		api.APISet,
	)
	return new(Injector), nil, nil
}
