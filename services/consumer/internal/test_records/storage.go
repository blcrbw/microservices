package test_records

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, record *TestRecord) error
	CreateMultiple(ctx context.Context, records []*TestRecord) error
	FindAll(ctx context.Context) (u []TestRecord, err error)
	FindOne(ctx context.Context, id string) (TestRecord, error)
	Update(ctx context.Context, user *TestRecord) error
	Delete(ctx context.Context, id string) error
}
