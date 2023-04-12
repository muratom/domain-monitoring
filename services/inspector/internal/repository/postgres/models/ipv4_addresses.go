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

// Ipv4Address is an object representing the database table.
type Ipv4Address struct {
	ID       int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	DomainID int    `boil:"domain_id" json:"domain_id" toml:"domain_id" yaml:"domain_id"`
	IP       string `boil:"ip" json:"ip" toml:"ip" yaml:"ip"`

	R *ipv4AddressR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L ipv4AddressL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var Ipv4AddressColumns = struct {
	ID       string
	DomainID string
	IP       string
}{
	ID:       "id",
	DomainID: "domain_id",
	IP:       "ip",
}

var Ipv4AddressTableColumns = struct {
	ID       string
	DomainID string
	IP       string
}{
	ID:       "ipv4_addresses.id",
	DomainID: "ipv4_addresses.domain_id",
	IP:       "ipv4_addresses.ip",
}

// Generated where

var Ipv4AddressWhere = struct {
	ID       whereHelperint
	DomainID whereHelperint
	IP       whereHelperstring
}{
	ID:       whereHelperint{field: "\"ipv4_addresses\".\"id\""},
	DomainID: whereHelperint{field: "\"ipv4_addresses\".\"domain_id\""},
	IP:       whereHelperstring{field: "\"ipv4_addresses\".\"ip\""},
}

// Ipv4AddressRels is where relationship names are stored.
var Ipv4AddressRels = struct {
	Domain string
}{
	Domain: "Domain",
}

// ipv4AddressR is where relationships are stored.
type ipv4AddressR struct {
	Domain *Domain `boil:"Domain" json:"Domain" toml:"Domain" yaml:"Domain"`
}

// NewStruct creates a new relationship struct
func (*ipv4AddressR) NewStruct() *ipv4AddressR {
	return &ipv4AddressR{}
}

func (r *ipv4AddressR) GetDomain() *Domain {
	if r == nil {
		return nil
	}
	return r.Domain
}

// ipv4AddressL is where Load methods for each relationship are stored.
type ipv4AddressL struct{}

var (
	ipv4AddressAllColumns            = []string{"id", "domain_id", "ip"}
	ipv4AddressColumnsWithoutDefault = []string{"domain_id", "ip"}
	ipv4AddressColumnsWithDefault    = []string{"id"}
	ipv4AddressPrimaryKeyColumns     = []string{"id"}
	ipv4AddressGeneratedColumns      = []string{"id"}
)

type (
	// Ipv4AddressSlice is an alias for a slice of pointers to Ipv4Address.
	// This should almost always be used instead of []Ipv4Address.
	Ipv4AddressSlice []*Ipv4Address
	// Ipv4AddressHook is the signature for custom Ipv4Address hook methods
	Ipv4AddressHook func(context.Context, boil.ContextExecutor, *Ipv4Address) error

	ipv4AddressQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	ipv4AddressType                 = reflect.TypeOf(&Ipv4Address{})
	ipv4AddressMapping              = queries.MakeStructMapping(ipv4AddressType)
	ipv4AddressPrimaryKeyMapping, _ = queries.BindMapping(ipv4AddressType, ipv4AddressMapping, ipv4AddressPrimaryKeyColumns)
	ipv4AddressInsertCacheMut       sync.RWMutex
	ipv4AddressInsertCache          = make(map[string]insertCache)
	ipv4AddressUpdateCacheMut       sync.RWMutex
	ipv4AddressUpdateCache          = make(map[string]updateCache)
	ipv4AddressUpsertCacheMut       sync.RWMutex
	ipv4AddressUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var ipv4AddressAfterSelectHooks []Ipv4AddressHook

var ipv4AddressBeforeInsertHooks []Ipv4AddressHook
var ipv4AddressAfterInsertHooks []Ipv4AddressHook

var ipv4AddressBeforeUpdateHooks []Ipv4AddressHook
var ipv4AddressAfterUpdateHooks []Ipv4AddressHook

var ipv4AddressBeforeDeleteHooks []Ipv4AddressHook
var ipv4AddressAfterDeleteHooks []Ipv4AddressHook

var ipv4AddressBeforeUpsertHooks []Ipv4AddressHook
var ipv4AddressAfterUpsertHooks []Ipv4AddressHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Ipv4Address) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Ipv4Address) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Ipv4Address) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Ipv4Address) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Ipv4Address) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Ipv4Address) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Ipv4Address) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Ipv4Address) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Ipv4Address) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ipv4AddressAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddIpv4AddressHook registers your hook function for all future operations.
func AddIpv4AddressHook(hookPoint boil.HookPoint, ipv4AddressHook Ipv4AddressHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		ipv4AddressAfterSelectHooks = append(ipv4AddressAfterSelectHooks, ipv4AddressHook)
	case boil.BeforeInsertHook:
		ipv4AddressBeforeInsertHooks = append(ipv4AddressBeforeInsertHooks, ipv4AddressHook)
	case boil.AfterInsertHook:
		ipv4AddressAfterInsertHooks = append(ipv4AddressAfterInsertHooks, ipv4AddressHook)
	case boil.BeforeUpdateHook:
		ipv4AddressBeforeUpdateHooks = append(ipv4AddressBeforeUpdateHooks, ipv4AddressHook)
	case boil.AfterUpdateHook:
		ipv4AddressAfterUpdateHooks = append(ipv4AddressAfterUpdateHooks, ipv4AddressHook)
	case boil.BeforeDeleteHook:
		ipv4AddressBeforeDeleteHooks = append(ipv4AddressBeforeDeleteHooks, ipv4AddressHook)
	case boil.AfterDeleteHook:
		ipv4AddressAfterDeleteHooks = append(ipv4AddressAfterDeleteHooks, ipv4AddressHook)
	case boil.BeforeUpsertHook:
		ipv4AddressBeforeUpsertHooks = append(ipv4AddressBeforeUpsertHooks, ipv4AddressHook)
	case boil.AfterUpsertHook:
		ipv4AddressAfterUpsertHooks = append(ipv4AddressAfterUpsertHooks, ipv4AddressHook)
	}
}

// One returns a single ipv4Address record from the query.
func (q ipv4AddressQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Ipv4Address, error) {
	o := &Ipv4Address{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for ipv4_addresses")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Ipv4Address records from the query.
func (q ipv4AddressQuery) All(ctx context.Context, exec boil.ContextExecutor) (Ipv4AddressSlice, error) {
	var o []*Ipv4Address

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Ipv4Address slice")
	}

	if len(ipv4AddressAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Ipv4Address records in the query.
func (q ipv4AddressQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count ipv4_addresses rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q ipv4AddressQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if ipv4_addresses exists")
	}

	return count > 0, nil
}

// Domain pointed to by the foreign key.
func (o *Ipv4Address) Domain(mods ...qm.QueryMod) domainQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.DomainID),
	}

	queryMods = append(queryMods, mods...)

	return Domains(queryMods...)
}

// LoadDomain allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (ipv4AddressL) LoadDomain(ctx context.Context, e boil.ContextExecutor, singular bool, maybeIpv4Address interface{}, mods queries.Applicator) error {
	var slice []*Ipv4Address
	var object *Ipv4Address

	if singular {
		var ok bool
		object, ok = maybeIpv4Address.(*Ipv4Address)
		if !ok {
			object = new(Ipv4Address)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeIpv4Address)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeIpv4Address))
			}
		}
	} else {
		s, ok := maybeIpv4Address.(*[]*Ipv4Address)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeIpv4Address)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeIpv4Address))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &ipv4AddressR{}
		}
		args = append(args, object.DomainID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &ipv4AddressR{}
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

	if len(domainAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
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
		foreign.R.Ipv4Addresses = append(foreign.R.Ipv4Addresses, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.DomainID == foreign.ID {
				local.R.Domain = foreign
				if foreign.R == nil {
					foreign.R = &domainR{}
				}
				foreign.R.Ipv4Addresses = append(foreign.R.Ipv4Addresses, local)
				break
			}
		}
	}

	return nil
}

// SetDomain of the ipv4Address to the related item.
// Sets o.R.Domain to related.
// Adds o to related.R.Ipv4Addresses.
func (o *Ipv4Address) SetDomain(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Domain) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"ipv4_addresses\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"domain_id"}),
		strmangle.WhereClause("\"", "\"", 2, ipv4AddressPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

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
		o.R = &ipv4AddressR{
			Domain: related,
		}
	} else {
		o.R.Domain = related
	}

	if related.R == nil {
		related.R = &domainR{
			Ipv4Addresses: Ipv4AddressSlice{o},
		}
	} else {
		related.R.Ipv4Addresses = append(related.R.Ipv4Addresses, o)
	}

	return nil
}

// Ipv4Addresses retrieves all the records using an executor.
func Ipv4Addresses(mods ...qm.QueryMod) ipv4AddressQuery {
	mods = append(mods, qm.From("\"ipv4_addresses\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"ipv4_addresses\".*"})
	}

	return ipv4AddressQuery{q}
}

// FindIpv4Address retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindIpv4Address(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Ipv4Address, error) {
	ipv4AddressObj := &Ipv4Address{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"ipv4_addresses\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, ipv4AddressObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from ipv4_addresses")
	}

	if err = ipv4AddressObj.doAfterSelectHooks(ctx, exec); err != nil {
		return ipv4AddressObj, err
	}

	return ipv4AddressObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Ipv4Address) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no ipv4_addresses provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(ipv4AddressColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	ipv4AddressInsertCacheMut.RLock()
	cache, cached := ipv4AddressInsertCache[key]
	ipv4AddressInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			ipv4AddressAllColumns,
			ipv4AddressColumnsWithDefault,
			ipv4AddressColumnsWithoutDefault,
			nzDefaults,
		)
		wl = strmangle.SetComplement(wl, ipv4AddressGeneratedColumns)

		cache.valueMapping, err = queries.BindMapping(ipv4AddressType, ipv4AddressMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(ipv4AddressType, ipv4AddressMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"ipv4_addresses\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"ipv4_addresses\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into ipv4_addresses")
	}

	if !cached {
		ipv4AddressInsertCacheMut.Lock()
		ipv4AddressInsertCache[key] = cache
		ipv4AddressInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Ipv4Address.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Ipv4Address) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	ipv4AddressUpdateCacheMut.RLock()
	cache, cached := ipv4AddressUpdateCache[key]
	ipv4AddressUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			ipv4AddressAllColumns,
			ipv4AddressPrimaryKeyColumns,
		)
		wl = strmangle.SetComplement(wl, ipv4AddressGeneratedColumns)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update ipv4_addresses, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"ipv4_addresses\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, ipv4AddressPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(ipv4AddressType, ipv4AddressMapping, append(wl, ipv4AddressPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update ipv4_addresses row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for ipv4_addresses")
	}

	if !cached {
		ipv4AddressUpdateCacheMut.Lock()
		ipv4AddressUpdateCache[key] = cache
		ipv4AddressUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q ipv4AddressQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for ipv4_addresses")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for ipv4_addresses")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o Ipv4AddressSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ipv4AddressPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"ipv4_addresses\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, ipv4AddressPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in ipv4Address slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all ipv4Address")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Ipv4Address) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no ipv4_addresses provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(ipv4AddressColumnsWithDefault, o)

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

	ipv4AddressUpsertCacheMut.RLock()
	cache, cached := ipv4AddressUpsertCache[key]
	ipv4AddressUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			ipv4AddressAllColumns,
			ipv4AddressColumnsWithDefault,
			ipv4AddressColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			ipv4AddressAllColumns,
			ipv4AddressPrimaryKeyColumns,
		)

		insert = strmangle.SetComplement(insert, ipv4AddressGeneratedColumns)
		update = strmangle.SetComplement(update, ipv4AddressGeneratedColumns)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert ipv4_addresses, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(ipv4AddressPrimaryKeyColumns))
			copy(conflict, ipv4AddressPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"ipv4_addresses\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(ipv4AddressType, ipv4AddressMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(ipv4AddressType, ipv4AddressMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert ipv4_addresses")
	}

	if !cached {
		ipv4AddressUpsertCacheMut.Lock()
		ipv4AddressUpsertCache[key] = cache
		ipv4AddressUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Ipv4Address record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Ipv4Address) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Ipv4Address provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), ipv4AddressPrimaryKeyMapping)
	sql := "DELETE FROM \"ipv4_addresses\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from ipv4_addresses")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for ipv4_addresses")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q ipv4AddressQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no ipv4AddressQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from ipv4_addresses")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for ipv4_addresses")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o Ipv4AddressSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(ipv4AddressBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ipv4AddressPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"ipv4_addresses\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, ipv4AddressPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from ipv4Address slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for ipv4_addresses")
	}

	if len(ipv4AddressAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Ipv4Address) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindIpv4Address(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *Ipv4AddressSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := Ipv4AddressSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ipv4AddressPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"ipv4_addresses\".* FROM \"ipv4_addresses\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, ipv4AddressPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in Ipv4AddressSlice")
	}

	*o = slice

	return nil
}

// Ipv4AddressExists checks if the Ipv4Address row exists.
func Ipv4AddressExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"ipv4_addresses\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if ipv4_addresses exists")
	}

	return exists, nil
}

// Exists checks if the Ipv4Address row exists.
func (o *Ipv4Address) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return Ipv4AddressExists(ctx, exec, o.ID)
}
