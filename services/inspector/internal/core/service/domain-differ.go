package service

import (
	"fmt"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/r3labs/diff"
)

type domainDiffer interface {
	Diff(a, b *entity.Domain) (entity.Changelog, error)
}

type libDomainDiffer struct{}

func (d *libDomainDiffer) Diff(a, b *entity.Domain) (entity.Changelog, error) {
	diffResult, err := diff.Diff(a, b)
	if err != nil {
		return nil, fmt.Errorf("finding diff of domain A (%+v) and B (%+v) was failed: %w", a, b, err)
	}
	changes := make([]entity.Change, len(diffResult))
	for i, res := range diffResult {
		changes[i] = entity.Change{
			Type: res.Type,
			Path: res.Path,
			From: res.From,
			To:   res.To,
		}
	}
	return changes, nil
}
