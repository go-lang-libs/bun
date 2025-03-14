package crud

import (
	"fmt"

	"github.com/uptrace/bun"
)

type KV[V any] struct {
	Key   string
	Value V
}

type Filter struct {
	OrInt64     []KV[int64]
	OrInt       []KV[int]
	OrString    []KV[string]
	AndInt      []KV[int]
	AndInt64    []KV[int64]
	AndString   []KV[string]
	WhereInt    []KV[int]
	WhereInt64  []KV[int64]
	WhereString []KV[string]
	Limit       int
	Offset      int
}

// applyFilter applies the filter conditions to the query.
func applyFilter(query *bun.SelectQuery, filter Filter) *bun.SelectQuery {
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	for _, kv := range filter.AndInt {
		query = query.Where(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.AndInt64 {
		query = query.Where(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.AndString {
		query = query.Where(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.OrInt {
		query = query.WhereOr(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.OrInt64 {
		query = query.WhereOr(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.OrString {
		query = query.WhereOr(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.WhereInt {
		query = query.Where(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.WhereInt64 {
		query = query.Where(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	for _, kv := range filter.WhereString {
		query = query.Where(fmt.Sprintf("%s = ?", kv.Key), kv.Value)
	}
	return query
}

// Filters
func WhereEmailFilter(email string) Filter {
	return Filter{WhereString: []KV[string]{{Key: "email", Value: email}}}
}

func WhereIdFilter(id int64) Filter {
	return Filter{WhereInt64: []KV[int64]{{Key: "id", Value: id}}}
}

func WhereSlugFilter(slug string) Filter {
	return Filter{WhereString: []KV[string]{{Key: "slug", Value: slug}}}
}

func WhereSlugOrId(slug string, id int64) Filter {
	return Filter{
		WhereString: []KV[string]{{Key: "slug", Value: slug}},
		WhereInt64:  []KV[int64]{{Key: "id", Value: id}},
	}
}
