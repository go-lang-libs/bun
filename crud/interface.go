package crud

import "context"

type Crudable[T any] interface {
	Create(ctx context.Context, a T) (T, error)
	Update(ctx context.Context, a T) (T, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, filter Filter) (T, error)
	List(ctx context.Context, filter Filter) ([]T, error)
	Exists(ctx context.Context, filter Filter) bool
}
