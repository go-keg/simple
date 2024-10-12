package dataloader

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/go-keg/keg/contrib/gql"
	"github.com/go-keg/simple/data/ent"
	"github.com/go-keg/simple/data/ent/permission"
	"github.com/go-keg/simple/data/ent/user"
	"github.com/graph-gophers/dataloader"
	"github.com/samber/lo"
)

type Loader struct {
	client *ent.Client
}

func (r Loader) permissionChildrenCount() gql.LoaderFunc {
	return func(ctx context.Context, keys dataloader.Keys) (map[dataloader.Key]any, error) {
		type item struct {
			ParentId int64 `json:"parent_id"`
			Count    int   `json:"count"`
		}
		var items []item
		err := r.client.Permission.Query().Where(permission.ParentIDIn(gql.ToInts(keys)...)).
			GroupBy(permission.FieldParentID).Aggregate(ent.Count()).Scan(ctx, &items)
		if err != nil {
			return nil, err
		}
		return lo.SliceToMap(items, func(item item) (dataloader.Key, any) {
			return gql.ToStringKey(item.ParentId), item.Count
		}), nil
	}
}

func (r Loader) userRoleCount() gql.LoaderFunc {
	type item struct {
		ID    int64 `json:"id"`
		Count int   `json:"count"`
	}
	return func(ctx context.Context, keys dataloader.Keys) (map[dataloader.Key]any, error) {
		var items []item
		err := r.client.User.Query().Where(user.IDIn(gql.ToInts(keys)...)).Modify(func(s *sql.Selector) {
			t1 := sql.Table("user_roles").As("t1")
			s.Select(s.C(user.FieldID), sql.As(sql.Count(t1.C("role_id")), "count")).
				From(s).LeftJoin(t1).On(s.C(user.FieldID), t1.C("user_id")).
				GroupBy(s.C(user.FieldID))
		}).Scan(ctx, &items)
		if err != nil {
			return nil, err
		}
		return lo.SliceToMap(items, func(item item) (dataloader.Key, any) {
			return gql.ToStringKey(item.ID), item.Count
		}), nil
	}
}
