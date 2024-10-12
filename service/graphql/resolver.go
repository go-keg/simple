package graphql

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-keg/keg/contrib/cache"
	"github.com/go-keg/keg/contrib/gql"
	"github.com/go-keg/simple/biz"
	"github.com/go-keg/simple/data/ent"
	"github.com/go-keg/simple/data/ent/permission"
	"github.com/go-keg/simple/data/ent/role"
	"github.com/go-keg/simple/data/ent/user"
	"github.com/go-keg/simple/server/auth"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/exp/slices"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	log            *log.Helper
	ent            *ent.Client
	accountUseCase *biz.AccountUseCase
	captcha        *base64Captcha.Captcha
}

// NewSchema creates a graphql executable schema.
func NewSchema(logger log.Logger, ent *ent.Client, accountUseCase *biz.AccountUseCase) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{
			log:            log.NewHelper(log.With(logger, "module", "service/graphql")),
			ent:            ent,
			accountUseCase: accountUseCase,
			captcha: base64Captcha.NewCaptcha(base64Captcha.NewDriverString(
				40,
				140,
				0,
				0,
				4,
				"1234567890abcdefghijklmnopqrktuvwxyz",
				&color.RGBA{},
				base64Captcha.DefaultEmbeddedFonts,
				nil), base64Captcha.DefaultMemStore),
		},
		Directives: DirectiveRoot{
			Disabled: func(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
				return nil, gql.ErrDisabled
			},
			Login: func(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
				if auth.GetUser(ctx) != nil {
					return next(ctx)
				}
				return nil, gql.ErrUnauthorized
			},
			HasPermission: func(ctx context.Context, obj any, next graphql.Resolver, key string) (res any, err error) {
				u := auth.GetUser(ctx)
				if u == nil {
					return nil, gql.ErrUnauthorized
				}
				if !u.IsAdmin {
					cacheKey := fmt.Sprintf("user:%d:permissions", u.ID)
					keys, err := cache.LocalRemember(cacheKey, time.Minute*2, func() ([]string, error) {
						return ent.Permission.Query().Where(
							permission.HasRolesWith(role.HasUsersWith(user.ID(u.ID))),
						).Select(permission.FieldKey).Strings(ctx)
					})
					if err != nil {
						return nil, err
					}
					if !slices.Contains(keys, key) {
						return res, gql.ErrNoPermission
					}
				}
				return next(ctx)
			},
		},
	})
}
