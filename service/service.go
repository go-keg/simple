package service

import (
	"github.com/go-keg/simple/service/graphql"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(graphql.NewSchema)
