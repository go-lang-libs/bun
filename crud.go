package crud

import (
	"context"
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

type Crudable[T any] interface {
	Add(ctx context.Context, a T) error
	Update(ctx context.Context, a T) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, filter Filter) (T, error)
	List(ctx context.Context, filter Filter) ([]T, error)
	Exists(ctx context.Context, filter Filter) bool
}

func New[T any](db *bun.DB) CRUD[T] {
	return CRUD[T]{db: db}
}

type CRUD[T any] struct {
	db *bun.DB
	// Add necessary fields like DB connection, logger, etc.
}

// Add inserts a new record using Bun's NewInsert.
func (c *CRUD[T]) Add(ctx context.Context, a T) error {
	_, err := c.db.NewInsert().Model(a).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert error: %w", err)
	}
	return nil
}

// Update modifies an existing record. Bun will use the model’s primary key.
func (c *CRUD[T]) Update(ctx context.Context, a T) error {
	_, err := c.db.NewUpdate().Model(a).Exec(ctx)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}
	return nil
}

// Delete removes a record by its id.
func (c *CRUD[T]) Delete(ctx context.Context, id int64) error {
	// create a zero value instance to infer the table/model
	var m T
	_, err := c.db.NewDelete().Model(&m).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}
	return nil
}

// List retrieves multiple records based on the provided filter.
func (c *CRUD[T]) List(ctx context.Context, filter Filter) ([]T, error) {
	var results []T
	query := c.db.NewSelect().Model(&results).Order("id DESC")
	query = applyFilter(query, filter)
	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list error: %w", err)
	}
	return results, nil
}

// Get retrieves a single record based on the provided filter.
func (c *CRUD[T]) Get(ctx context.Context, filter Filter) (T, error) {
	var result T
	query := c.db.NewSelect().Model(&result)
	query = applyFilter(query, filter)
	err := query.Scan(ctx)
	if err != nil {
		return result, fmt.Errorf("get error: %w", err)
	}
	return result, nil
}

func (s *CRUD[T]) Exists(ctx context.Context, filter Filter) bool {
	var (
		result T
		id     int64
	)
	query := applyFilter(s.db.NewSelect().Model(&result), filter)
	err := query.Scan(ctx, &id)
	if err != nil {
		return false
	}
	return id > 0
}

// GetBy retrieves a single record that matches the given field and value.
func (c *CRUD[T]) GetBy(ctx context.Context, field string, value any) (T, error) {
	var result T
	condition := fmt.Sprintf("%s = ?", field)
	err := c.db.NewSelect().Model(&result).Where(condition, value).Scan(ctx)
	if err != nil {
		return result, fmt.Errorf("getBy error: %w", err)
	}
	return result, nil
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
