package service

import (
	"fmt"
	"strconv"

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

	changelog := make(entity.Changelog, len(diffResult))
	for _, diffRes := range diffResult {
		newPath := filter(diffRes.Path, func(s string) bool {
			_, err := strconv.Atoi(s)
			return err != nil
		})

		fieldType, path := getFieldTypeAndPath(newPath)
		change := entity.Change{
			OperationType: mapOperationType(diffRes.Type),
			FieldType:     fieldType,
			Path:          path,
			From:          diffRes.From,
			To:            diffRes.To,
		}
		changelog = append(changelog, change)
	}

	return changelog, nil
}

func mapOperationType(opType string) entity.OperationType {
	switch opType {
	case diff.CREATE:
		return entity.CREATE
	case diff.UPDATE:
		return entity.UPDATE
	case diff.DELETE:
		return entity.DELETE
	}
	panic(fmt.Sprintf("failed to map operation type %v", opType))
}

func getFieldTypeAndPath(path []string) (entity.FieldType, []string) {
	switch path[0] {
	case "FQDN":
		return entity.FQDN, []string{}
	case "DNS":
		return entity.DNS, path[1:]
	case "WHOIS":
		return entity.WHOIS, path[1:]
	}
	panic(fmt.Sprintf("failed to get field type and path %v", path))
}

func filter[T any](ss []T, isPass func(T) bool) (ret []T) {
	for _, s := range ss {
		if isPass(s) {
			ret = append(ret, s)
		}
	}
	return
}
