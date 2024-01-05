package domain

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/changelog"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
)

type Differ interface {
	Diff(a, b *domain.Domain) (changelog.Changelog, error)
}
