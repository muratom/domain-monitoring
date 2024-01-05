package inspector

import (
	"context"
)

type retriever interface {
	RetrieveRottenDomainsFQDN(ctx context.Context) ([]string, error)
}
