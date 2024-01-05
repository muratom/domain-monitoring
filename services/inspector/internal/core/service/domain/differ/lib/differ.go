package lib

import (
	"fmt"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/changelog"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
	"github.com/r3labs/diff"
	"strconv"
)

type Differ struct{}

func (d *Differ) Diff(a, b *domain.Domain) (changelog.Changelog, error) {
	diffResult, err := diff.Diff(a, b)
	if err != nil {
		return nil, fmt.Errorf("finding diff of domain A (%+v) and B (%+v) was failed: %w", a, b, err)
	}

	changeLog := make(changelog.Changelog, 0, len(diffResult))
	for _, diffRes := range diffResult {
		newPath := filter(diffRes.Path, func(s string) bool {
			_, err := strconv.Atoi(s)
			return err != nil
		})

		fieldType, path := getFieldTypeAndPath(newPath)
		change := changelog.Change{
			OperationType: mapOperationType(diffRes.Type),
			FieldType:     fieldType,
			Path:          path,
			From:          diffRes.From,
			To:            diffRes.To,
		}
		changeLog = append(changeLog, change)
	}

	return changeLog, nil
}

func mapOperationType(opType string) changelog.OperationType {
	switch opType {
	case diff.CREATE:
		return changelog.CREATE
	case diff.UPDATE:
		return changelog.UPDATE
	case diff.DELETE:
		return changelog.DELETE
	}
	panic(fmt.Sprintf("failed to map operation type %v", opType))
}

func getFieldTypeAndPath(path []string) (changelog.FieldType, []string) {
	switch path[0] {
	case "FQDN":
		return changelog.FQDN, []string{}
	case "DNS":
		return changelog.DNS, path[1:]
	case "WHOIS":
		return changelog.WHOIS, path[1:]
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
