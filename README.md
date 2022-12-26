# ges

simple elasticsearch orm

## build condition

- Term(field string, value interface{}) Filter
- Terms(field string, values ...interface{}) Filter
- Between(field string, start, end int64) Filter
- Gt(field string, value int64) Filter
- Gte(field string, value int64) Filter
- Lt(field string, value int64) Filter
- Lte(field string, value int64) Filter
- Wildcard(field string, value string) Filter
- WildcardSuffix(field string, value string) Filter

## build aggregator condition 
### agg 
```go
type Agg interface {
	Name(string) Agg
	DateHistogram(field, interval, format, offset, timeZone string) Agg
	Distinct(field string, number int64) Agg
	Filter(agg Agg, filter ...AggFilter) Agg
 	Result() (string, interface{})
}
```

### aggregator filter 
```go
type AggFilter interface {
	Name(string) AggFilter
	Terms(field string, val ...interface{}) AggFilter
	TermsArray(field string, val interface{}) AggFilter

	Result() (string, map[string]interface{})
}
```

## elastic search query 


### index 
```go
Exists(ctx context.Context) (bool, error)
	Create(ctx context.Context, mapping IndexMeta) error
	List(ctx context.Context) ([]error, error)
	Mapping(ctx context.Context) (map[string]indexMetaResp, error)

```


### execute 
```go
	IndexName(name string) Client
	Index() Index

	Not(filters ...Filter) Client
	Where(filters ...Filter) Client
	SQLWhere(query string, args ...interface{}) Client
	Or(filters ...Filter) Client
	OrderBy(field string, isDesc bool) Client
	Size(uint64) Client
	Agg(aggs ...Agg) Client
	Start(uint64) Client
	Limit(uint64, uint64) Client
	Fields(...string) Client
	Search(ctx context.Context, result interface{}) (uint64, error)
	GetById(ctx context.Context, id string, result interface{}) error
	RawSQL(ctx context.Context, closer io.ReadCloser, result interface{}) (uint64, error)
	Count(ctx context.Context) (uint64, error)
	Save(ctx context.Context, data ...interface{}) error
	USave(ctx context.Context, docs ...Document) error
	UpdateById(ctx context.Context, id string, data interface{}) error
	// TODO map[string]interface{} to interface api
	MUpdateById(ctx context.Context, docs ...Document) error
	// TODO map[string]interface{} to interface api
	MUpsertById(ctx context.Context, docs ...Document) error
	UpsertById(ctx context.Context, id string, doc interface{}) error
	// Delete delete_by_query
	Delete(ctx context.Context) error
	DeleteById(ctx context.Context, ids ...string) error
```
