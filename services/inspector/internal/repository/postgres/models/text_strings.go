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

// TextString is an object representing the database table.
type TextString struct {
	ID       int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	DomainID int    `boil:"domain_id" json:"domain_id" toml:"domain_id" yaml:"domain_id"`
	Text     string `boil:"text" json:"text" toml:"text" yaml:"text"`

	R *textStringR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L textStringL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TextStringColumns = struct {
	ID       string
	DomainID string
	Text     string
}{
	ID:       "id",
	DomainID: "domain_id",
	Text:     "text",
}

var TextStringTableColumns = struct {
	ID       string
	DomainID string
	Text     string
}{
	ID:       "text_strings.id",
	DomainID: "text_strings.domain_id",
	Text:     "text_strings.text",
}

// Generated where

var TextStringWhere = struct {
	ID       whereHelperint
	DomainID whereHelperint
	Text     whereHelperstring
}{
	ID:       whereHelperint{field: "\"text_strings\".\"id\""},
	DomainID: whereHelperint{field: "\"text_strings\".\"domain_id\""},
	Text:     whereHelperstring{field: "\"text_strings\".\"text\""},
}

// TextStringRels is where relationship names are stored.
var TextStringRels = struct {
	Domain string
}{
	Domain: "Domain",
}

// textStringR is where relationships are stored.
type textStringR struct {
	Domain *Domain `boil:"Domain" json:"Domain" toml:"Domain" yaml:"Domain"`
}

// NewStruct creates a new relationship struct
func (*textStringR) NewStruct() *textStringR {
	return &textStringR{}
}

func (r *textStringR) GetDomain() *Domain {
	if r == nil {
		return nil
	}
	return r.Domain
}

// textStringL is where Load methods for each relationship are stored.
type textStringL struct{}

var (
	textStringAllColumns            = []string{"id", "domain_id", "text"}
	textStringColumnsWithoutDefault = []string{"domain_id", "text"}
	textStringColumnsWithDefault    = []string{"id"}
	textStringPrimaryKeyColumns     = []string{"domain_id", "text"}
	textStringGeneratedColumns      = []string{"id"}
)

type (
	// TextStringSlice is an alias for a slice of pointers to TextString.
	// This should almost always be used instead of []TextString.
	TextStringSlice []*TextString

	textStringQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	textStringType                 = reflect.TypeOf(&TextString{})
	textStringMapping              = queries.MakeStructMapping(textStringType)
	textStringPrimaryKeyMapping, _ = queries.BindMapping(textStringType, textStringMapping, textStringPrimaryKeyColumns)
	textStringInsertCacheMut       sync.RWMutex
	textStringInsertCache          = make(map[string]insertCache)
	textStringUpdateCacheMut       sync.RWMutex
	textStringUpdateCache          = make(map[string]updateCache)
	textStringUpsertCacheMut       sync.RWMutex
	textStringUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single textString record from the query.
func (q textStringQuery) One(ctx context.Context, exec boil.ContextExecutor) (*TextString, error) {
	o := &TextString{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for text_strings")
	}

	return o, nil
}

// All returns all TextString records from the query.
func (q textStringQuery) All(ctx context.Context, exec boil.ContextExecutor) (TextStringSlice, error) {
	var o []*TextString

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to TextString slice")
	}

	return o, nil
}

// Count returns the count of all TextString records in the query.
func (q textStringQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count text_strings rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q textStringQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if text_strings exists")
	}

	return count > 0, nil
}

// Domain pointed to by the foreign key.
func (o *TextString) Domain(mods ...qm.QueryMod) domainQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.DomainID),
	}

	queryMods = append(queryMods, mods...)

	return Domains(queryMods...)
}

// LoadDomain allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (textStringL) LoadDomain(ctx context.Context, e boil.ContextExecutor, singular bool, maybeTextString interface{}, mods queries.Applicator) error {
	var slice []*TextString
	var object *TextString

	if singular {
		var ok bool
		object, ok = maybeTextString.(*TextString)
		if !ok {
			object = new(TextString)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeTextString)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeTextString))
			}
		}
	} else {
		s, ok := maybeTextString.(*[]*TextString)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeTextString)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeTextString))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &textStringR{}
		}
		args = append(args, object.DomainID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &textStringR{}
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
		foreign.R.TextStrings = append(foreign.R.TextStrings, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.DomainID == foreign.ID {
				local.R.Domain = foreign
				if foreign.R == nil {
					foreign.R = &domainR{}
				}
				foreign.R.TextStrings = append(foreign.R.TextStrings, local)
				break
			}
		}
	}

	return nil
}

// SetDomain of the textString to the related item.
// Sets o.R.Domain to related.
// Adds o to related.R.TextStrings.
func (o *TextString) SetDomain(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Domain) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"text_strings\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"domain_id"}),
		strmangle.WhereClause("\"", "\"", 2, textStringPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.DomainID, o.Text}

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
		o.R = &textStringR{
			Domain: related,
		}
	} else {
		o.R.Domain = related
	}

	if related.R == nil {
		related.R = &domainR{
			TextStrings: TextStringSlice{o},
		}
	} else {
		related.R.TextStrings = append(related.R.TextStrings, o)
	}

	return nil
}

// TextStrings retrieves all the records using an executor.
func TextStrings(mods ...qm.QueryMod) textStringQuery {
	mods = append(mods, qm.From("\"text_strings\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"text_strings\".*"})
	}

	return textStringQuery{q}
}

// FindTextString retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTextString(ctx context.Context, exec boil.ContextExecutor, domainID int, text string, selectCols ...string) (*TextString, error) {
	textStringObj := &TextString{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"text_strings\" where \"domain_id\"=$1 AND \"text\"=$2", sel,
	)

	q := queries.Raw(query, domainID, text)

	err := q.Bind(ctx, exec, textStringObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from text_strings")
	}

	return textStringObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *TextString) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no text_strings provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(textStringColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	textStringInsertCacheMut.RLock()
	cache, cached := textStringInsertCache[key]
	textStringInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			textStringAllColumns,
			textStringColumnsWithDefault,
			textStringColumnsWithoutDefault,
			nzDefaults,
		)
		wl = strmangle.SetComplement(wl, textStringGeneratedColumns)

		cache.valueMapping, err = queries.BindMapping(textStringType, textStringMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(textStringType, textStringMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"text_strings\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"text_strings\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into text_strings")
	}

	if !cached {
		textStringInsertCacheMut.Lock()
		textStringInsertCache[key] = cache
		textStringInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the TextString.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *TextString) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	textStringUpdateCacheMut.RLock()
	cache, cached := textStringUpdateCache[key]
	textStringUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			textStringAllColumns,
			textStringPrimaryKeyColumns,
		)
		wl = strmangle.SetComplement(wl, textStringGeneratedColumns)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update text_strings, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"text_strings\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, textStringPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(textStringType, textStringMapping, append(wl, textStringPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update text_strings row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for text_strings")
	}

	if !cached {
		textStringUpdateCacheMut.Lock()
		textStringUpdateCache[key] = cache
		textStringUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q textStringQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for text_strings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for text_strings")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TextStringSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), textStringPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"text_strings\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, textStringPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in textString slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all textString")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *TextString) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no text_strings provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(textStringColumnsWithDefault, o)

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

	textStringUpsertCacheMut.RLock()
	cache, cached := textStringUpsertCache[key]
	textStringUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			textStringAllColumns,
			textStringColumnsWithDefault,
			textStringColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			textStringAllColumns,
			textStringPrimaryKeyColumns,
		)

		insert = strmangle.SetComplement(insert, textStringGeneratedColumns)
		update = strmangle.SetComplement(update, textStringGeneratedColumns)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert text_strings, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(textStringPrimaryKeyColumns))
			copy(conflict, textStringPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"text_strings\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(textStringType, textStringMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(textStringType, textStringMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert text_strings")
	}

	if !cached {
		textStringUpsertCacheMut.Lock()
		textStringUpsertCache[key] = cache
		textStringUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single TextString record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *TextString) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no TextString provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), textStringPrimaryKeyMapping)
	sql := "DELETE FROM \"text_strings\" WHERE \"domain_id\"=$1 AND \"text\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from text_strings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for text_strings")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q textStringQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no textStringQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from text_strings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for text_strings")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TextStringSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), textStringPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"text_strings\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, textStringPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from textString slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for text_strings")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *TextString) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindTextString(ctx, exec, o.DomainID, o.Text)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TextStringSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TextStringSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), textStringPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"text_strings\".* FROM \"text_strings\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, textStringPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in TextStringSlice")
	}

	*o = slice

	return nil
}

// TextStringExists checks if the TextString row exists.
func TextStringExists(ctx context.Context, exec boil.ContextExecutor, domainID int, text string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"text_strings\" where \"domain_id\"=$1 AND \"text\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, domainID, text)
	}
	row := exec.QueryRowContext(ctx, sql, domainID, text)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if text_strings exists")
	}

	return exists, nil
}

// Exists checks if the TextString row exists.
func (o *TextString) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return TextStringExists(ctx, exec, o.DomainID, o.Text)
}
