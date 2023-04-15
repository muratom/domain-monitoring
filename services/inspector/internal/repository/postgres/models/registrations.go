// Code generated by SQLBoiler 4.14.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Registration is an object representing the database table.
type Registration struct {
	ID        int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	DomainID  int       `boil:"domain_id" json:"domain_id" toml:"domain_id" yaml:"domain_id"`
	Created   time.Time `boil:"created" json:"created" toml:"created" yaml:"created"`
	PaidTill  time.Time `boil:"paid_till" json:"paid_till" toml:"paid_till" yaml:"paid_till"`
	Registrar string    `boil:"registrar" json:"registrar" toml:"registrar" yaml:"registrar"`

	R *registrationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L registrationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var RegistrationColumns = struct {
	ID        string
	DomainID  string
	Created   string
	PaidTill  string
	Registrar string
}{
	ID:        "id",
	DomainID:  "domain_id",
	Created:   "created",
	PaidTill:  "paid_till",
	Registrar: "registrar",
}

var RegistrationTableColumns = struct {
	ID        string
	DomainID  string
	Created   string
	PaidTill  string
	Registrar string
}{
	ID:        "registrations.id",
	DomainID:  "registrations.domain_id",
	Created:   "registrations.created",
	PaidTill:  "registrations.paid_till",
	Registrar: "registrations.registrar",
}

// Generated where

var RegistrationWhere = struct {
	ID        whereHelperint
	DomainID  whereHelperint
	Created   whereHelpertime_Time
	PaidTill  whereHelpertime_Time
	Registrar whereHelperstring
}{
	ID:        whereHelperint{field: "\"registrations\".\"id\""},
	DomainID:  whereHelperint{field: "\"registrations\".\"domain_id\""},
	Created:   whereHelpertime_Time{field: "\"registrations\".\"created\""},
	PaidTill:  whereHelpertime_Time{field: "\"registrations\".\"paid_till\""},
	Registrar: whereHelperstring{field: "\"registrations\".\"registrar\""},
}

// RegistrationRels is where relationship names are stored.
var RegistrationRels = struct {
	Domain string
}{
	Domain: "Domain",
}

// registrationR is where relationships are stored.
type registrationR struct {
	Domain *Domain `boil:"Domain" json:"Domain" toml:"Domain" yaml:"Domain"`
}

// NewStruct creates a new relationship struct
func (*registrationR) NewStruct() *registrationR {
	return &registrationR{}
}

func (r *registrationR) GetDomain() *Domain {
	if r == nil {
		return nil
	}
	return r.Domain
}

// registrationL is where Load methods for each relationship are stored.
type registrationL struct{}

var (
	registrationAllColumns            = []string{"id", "domain_id", "created", "paid_till", "registrar"}
	registrationColumnsWithoutDefault = []string{"domain_id", "created", "paid_till", "registrar"}
	registrationColumnsWithDefault    = []string{"id"}
	registrationPrimaryKeyColumns     = []string{"domain_id", "created"}
	registrationGeneratedColumns      = []string{"id"}
)

type (
	// RegistrationSlice is an alias for a slice of pointers to Registration.
	// This should almost always be used instead of []Registration.
	RegistrationSlice []*Registration

	registrationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	registrationType                 = reflect.TypeOf(&Registration{})
	registrationMapping              = queries.MakeStructMapping(registrationType)
	registrationPrimaryKeyMapping, _ = queries.BindMapping(registrationType, registrationMapping, registrationPrimaryKeyColumns)
	registrationInsertCacheMut       sync.RWMutex
	registrationInsertCache          = make(map[string]insertCache)
	registrationUpdateCacheMut       sync.RWMutex
	registrationUpdateCache          = make(map[string]updateCache)
	registrationUpsertCacheMut       sync.RWMutex
	registrationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single registration record from the query.
func (q registrationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Registration, error) {
	o := &Registration{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for registrations")
	}

	return o, nil
}

// All returns all Registration records from the query.
func (q registrationQuery) All(ctx context.Context, exec boil.ContextExecutor) (RegistrationSlice, error) {
	var o []*Registration

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Registration slice")
	}

	return o, nil
}

// Count returns the count of all Registration records in the query.
func (q registrationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count registrations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q registrationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if registrations exists")
	}

	return count > 0, nil
}

// Domain pointed to by the foreign key.
func (o *Registration) Domain(mods ...qm.QueryMod) domainQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.DomainID),
	}

	queryMods = append(queryMods, mods...)

	return Domains(queryMods...)
}

// LoadDomain allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (registrationL) LoadDomain(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRegistration interface{}, mods queries.Applicator) error {
	var slice []*Registration
	var object *Registration

	if singular {
		var ok bool
		object, ok = maybeRegistration.(*Registration)
		if !ok {
			object = new(Registration)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeRegistration)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeRegistration))
			}
		}
	} else {
		s, ok := maybeRegistration.(*[]*Registration)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeRegistration)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeRegistration))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &registrationR{}
		}
		args = append(args, object.DomainID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &registrationR{}
			}

			for _, a := range args {
				if a == obj.DomainID {
					continue Outer
				}
			}

			args = append(args, obj.DomainID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`domains`),
		qm.WhereIn(`domains.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Domain")
	}

	var resultSlice []*Domain
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Domain")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for domains")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for domains")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Domain = foreign
		if foreign.R == nil {
			foreign.R = &domainR{}
		}
		foreign.R.Registrations = append(foreign.R.Registrations, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.DomainID == foreign.ID {
				local.R.Domain = foreign
				if foreign.R == nil {
					foreign.R = &domainR{}
				}
				foreign.R.Registrations = append(foreign.R.Registrations, local)
				break
			}
		}
	}

	return nil
}

// SetDomain of the registration to the related item.
// Sets o.R.Domain to related.
// Adds o to related.R.Registrations.
func (o *Registration) SetDomain(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Domain) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"registrations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"domain_id"}),
		strmangle.WhereClause("\"", "\"", 2, registrationPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.DomainID, o.Created}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.DomainID = related.ID
	if o.R == nil {
		o.R = &registrationR{
			Domain: related,
		}
	} else {
		o.R.Domain = related
	}

	if related.R == nil {
		related.R = &domainR{
			Registrations: RegistrationSlice{o},
		}
	} else {
		related.R.Registrations = append(related.R.Registrations, o)
	}

	return nil
}

// Registrations retrieves all the records using an executor.
func Registrations(mods ...qm.QueryMod) registrationQuery {
	mods = append(mods, qm.From("\"registrations\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"registrations\".*"})
	}

	return registrationQuery{q}
}

// FindRegistration retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRegistration(ctx context.Context, exec boil.ContextExecutor, domainID int, created time.Time, selectCols ...string) (*Registration, error) {
	registrationObj := &Registration{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"registrations\" where \"domain_id\"=$1 AND \"created\"=$2", sel,
	)

	q := queries.Raw(query, domainID, created)

	err := q.Bind(ctx, exec, registrationObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from registrations")
	}

	return registrationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Registration) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no registrations provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(registrationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	registrationInsertCacheMut.RLock()
	cache, cached := registrationInsertCache[key]
	registrationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			registrationAllColumns,
			registrationColumnsWithDefault,
			registrationColumnsWithoutDefault,
			nzDefaults,
		)
		wl = strmangle.SetComplement(wl, registrationGeneratedColumns)

		cache.valueMapping, err = queries.BindMapping(registrationType, registrationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(registrationType, registrationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"registrations\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"registrations\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into registrations")
	}

	if !cached {
		registrationInsertCacheMut.Lock()
		registrationInsertCache[key] = cache
		registrationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Registration.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Registration) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	registrationUpdateCacheMut.RLock()
	cache, cached := registrationUpdateCache[key]
	registrationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			registrationAllColumns,
			registrationPrimaryKeyColumns,
		)
		wl = strmangle.SetComplement(wl, registrationGeneratedColumns)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update registrations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"registrations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, registrationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(registrationType, registrationMapping, append(wl, registrationPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update registrations row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for registrations")
	}

	if !cached {
		registrationUpdateCacheMut.Lock()
		registrationUpdateCache[key] = cache
		registrationUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q registrationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for registrations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for registrations")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RegistrationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), registrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"registrations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, registrationPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in registration slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all registration")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Registration) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no registrations provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(registrationColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	registrationUpsertCacheMut.RLock()
	cache, cached := registrationUpsertCache[key]
	registrationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			registrationAllColumns,
			registrationColumnsWithDefault,
			registrationColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			registrationAllColumns,
			registrationPrimaryKeyColumns,
		)

		insert = strmangle.SetComplement(insert, registrationGeneratedColumns)
		update = strmangle.SetComplement(update, registrationGeneratedColumns)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert registrations, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(registrationPrimaryKeyColumns))
			copy(conflict, registrationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"registrations\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(registrationType, registrationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(registrationType, registrationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert registrations")
	}

	if !cached {
		registrationUpsertCacheMut.Lock()
		registrationUpsertCache[key] = cache
		registrationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Registration record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Registration) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Registration provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), registrationPrimaryKeyMapping)
	sql := "DELETE FROM \"registrations\" WHERE \"domain_id\"=$1 AND \"created\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from registrations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for registrations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q registrationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no registrationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from registrations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for registrations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RegistrationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), registrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"registrations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, registrationPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from registration slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for registrations")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Registration) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindRegistration(ctx, exec, o.DomainID, o.Created)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RegistrationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := RegistrationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), registrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"registrations\".* FROM \"registrations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, registrationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RegistrationSlice")
	}

	*o = slice

	return nil
}

// RegistrationExists checks if the Registration row exists.
func RegistrationExists(ctx context.Context, exec boil.ContextExecutor, domainID int, created time.Time) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"registrations\" where \"domain_id\"=$1 AND \"created\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, domainID, created)
	}
	row := exec.QueryRowContext(ctx, sql, domainID, created)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if registrations exists")
	}

	return exists, nil
}

// Exists checks if the Registration row exists.
func (o *Registration) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return RegistrationExists(ctx, exec, o.DomainID, o.Created)
}
