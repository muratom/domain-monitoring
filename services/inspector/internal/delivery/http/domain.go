package http

import (
	"context"
)

type domainService interface {
	AddDomain(ctx context.Context, addParams any)
	UpdateDomain(ctx context.Context, updateParams any)
	DeleteDomain(ctx context.Context, deleteParams any)
}
