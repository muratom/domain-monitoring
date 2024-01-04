package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/changelog"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

type ChangelogRepository struct {
	Conn *sql.DB
}

func NewChangelogRepository(dbConnection *sql.DB) *ChangelogRepository {
	return &ChangelogRepository{
		Conn: dbConnection,
	}
}

func (r *ChangelogRepository) Store(ctx context.Context, fqdn string, changelog changelog.Changelog) error {
	domainEntry, err := prepareDomainEntry(ctx, fqdn, r.Conn)
	if err != nil {
		return fmt.Errorf("failed to fetch data from DB for FQDN (%s): %w", fqdn, err)
	}

	rawChangelog, err := json.Marshal(changelog)
	if err != nil {
		return fmt.Errorf("error when making raw changelog: %w", err)
	}
	changelogEntry := &models.Changelog{
		CreatedAt: time.Now(),
		Changes:   rawChangelog,
	}

	return domainEntry.AddChangelogs(ctx, r.Conn, true, changelogEntry)
}

func (r *ChangelogRepository) Retrieve(ctx context.Context, fqdn string) ([]changelog.Changelog, error) {
	const (
		resultCapacity = 10
	)

	domainEntry, err := models.Domains(
		models.DomainWhere.FQDN.EQ(fqdn),
		qm.Load(models.DomainRels.Changelogs),
	).One(ctx, r.Conn)
	if err == sql.ErrNoRows {
		return []changelog.Changelog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get changelog for FQDN (%v): %w", fqdn, err)
	}

	result := make([]changelog.Changelog, 0, resultCapacity)
	for _, changelogEntry := range domainEntry.R.Changelogs {
		if changelogEntry != nil {
			var changes changelog.Changelog
			err := changelogEntry.Changes.Unmarshal(&changes)
			if err != nil {
				return nil, fmt.Errorf("unable to unmarshal changelog from DB: %w", err)
			}
			result = append(result, changes)
		}
	}

	return result, nil
}
