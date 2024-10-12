//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-keg/simple/biz"
	"github.com/go-keg/simple/conf"
	"github.com/go-keg/simple/data"
	"github.com/go-keg/simple/job"
	"github.com/go-keg/simple/schedule"
	"github.com/go-keg/simple/server"
	"github.com/go-keg/simple/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func initApp(log.Logger, *conf.Config) (*kratos.App, func(), error) {
	panic(wire.Build(
		biz.ProviderSet,
		data.ProviderSet,
		job.ProviderSet,
		schedule.ProviderSet,
		server.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
