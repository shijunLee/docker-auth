package storage

import (
	"context"
	"time"
)

type NamespaceHandler interface {
	List(ctx context.Context) ([]Namespace,error)
	GetByID(ctx context.Context,id string) (Namespace,error)
	Create(ctx context.Context,ns Namespace) (Namespace,error)
	Update(ctx context.Context,ns Namespace)(Namespace,error)
	Delete(ctx context.Context,id string) error
}

type Namespace interface {
	GetName(ctx context.Context) string
	GetCreateTime(ctx context.Context) time.Time
	GetNamespaceID(ctx context.Context) string
	GetInfo(ctx context.Context) []map[string]interface{}
}
