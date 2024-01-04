package service

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/changelog"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
)

type DomainDiffer interface {
	Diff(a, b *domain.Domain) (changelog.Changelog, error)
}
