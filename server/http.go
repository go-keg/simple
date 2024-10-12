package server

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-keg/simple/conf"
	"github.com/go-keg/simple/data/ent"
	"github.com/go-keg/simple/server/auth"
	"github.com/go-keg/simple/service/graphql/dataloader"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(
	cfg *conf.Config,
	logger log.Logger,
	client *ent.Client,
	schema graphql.ExecutableSchema,
) *http.Server {
	srv := http.NewServer(cfg.Server.Http.HttpOptions(logger)...)
	// graphql
	gqlSrv := handler.NewDefaultServer(schema)
	loader := dataloader.NewDataLoader(client)
	srv.Handle("/query", auth.Middleware(cfg.Key, client, dataloader.Middleware(loader, gqlSrv)))
	srv.HandleFunc("/graphql-ui", playground.Handler("Admin", "/query"))
	return srv
}
