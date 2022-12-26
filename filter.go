package ges

import "encoding/json"

/***************************
    @author: tiansheng.ren
    @date: 2022/11/3
    @desc:

***************************/

type term struct {
	name  string
	value interface{}
}

func (t term) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]interface{}{
		"term": {
			t.name: t.value,
		},
	})
}

type terms struct {
	name   string
	values []interface{}
}

func (t terms) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string][]interface{}{
		"terms": {
			t.name: t.values,
		},
	})
}

type between struct {
	name string `json:"-"`
	betweenValue
}
type betweenValue struct {
	Lte *int64 `json:"lte,omitempty"`
	Lt  *int64 `json:"lt,omitempty"`
	Gte *int64 `json:"gte,omitempty"`
	Gt  *int64 `json:"gt,omitempty"`
}

func (b between) MarshalJSON() ([]byte, error) {
	// 注意不能直接marshal b，回递归交通

	return json.Marshal(map[string]map[string]interface{}{
		"range": {
			b.name: b.betweenValue,
		},
	})
}

type wildCard struct {
	name string
	wildCardValue
}

type wildCardValue struct {
	Wildcard string  `json:"wildcard"`
	Boost    float64 `json:"boost"`
}

func (wcv wildCard) MarshalJSON() ([]byte, error) {
	// 注意不能直接marshal b，回递归交通

	return json.Marshal(map[string]map[string]interface{}{
		"wildcard": {
			wcv.name: wcv.wildCardValue,
		},
	})
}

type esCondition struct {
	Query esConditionQuery       `json:"query"`
	Agg   map[string]interface{} `json:"aggs,omitempty"`
}

type esConditionSortOrder struct {
	Order string `json:"order"`
}

type esConditionQuery struct {
	Bool esQueryBool `json:"bool"`
}

func (b esConditionQuery) MarshalJSON() ([]byte, error) {

	result, err := json.Marshal(b.Bool)
	if err != nil {
		return nil, err
	}
	if len(result) > 2 {
		return []byte(`{"bool":` + string(result) + `}`), nil
	}
	return []byte(`{"match_all": {}}`), nil
}

type esQueryBool struct {
	Must   []interface{} `json:"must,omitempty"`
	Not    []interface{} `json:"must_not,omitempty"`
	Should []interface{} `json:"should,omitempty"`
	Match  []interface{} `json:"match,omitempty"`
}

type filter struct {
	condition []interface{}
}

func NewFilter() Filter {
	return filter{}
}

func Term(field string, value interface{}) Filter { return filter{}.Term(field, value) }
func Terms(field string, values ...interface{}) Filter {
	return filter{}.Terms(field, values...)
}

func Between(field string, start, end int64) Filter {
	return filter{}.Between(field, start, end)
}
func Gt(field string, value int64) Filter              { return filter{}.Gt(field, value) }
func Gte(field string, value int64) Filter             { return filter{}.Gte(field, value) }
func Lt(field string, value int64) Filter              { return filter{}.Lt(field, value) }
func Lte(field string, value int64) Filter             { return filter{}.Lte(field, value) }
func Wildcard(field string, value string) Filter       { return filter{}.Wildcard(field, value) }
func WildcardSuffix(field string, value string) Filter { return filter{}.WildcardSuffix(field, value) }

func (f filter) Term(field string, value interface{}) Filter {
	f.condition = append(f.condition, term{field, value})
	return f
}

func (f filter) Terms(field string, values ...interface{}) Filter {
	f.condition = append(f.condition, terms{field, values})
	return f
}

func (f filter) Between(field string, start, end int64) Filter {
	b := between{name: field}
	b.Gte, b.Lte = &start, &end
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Gt(field string, value int64) Filter {
	b := between{name: field}
	b.Gt = &value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Gte(field string, value int64) Filter {
	b := between{name: field}
	b.Gte = &value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Lt(field string, value int64) Filter {
	b := between{name: field}
	b.Lt = &value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Lte(field string, value int64) Filter {
	b := between{name: field}
	b.Lte = &value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Wildcard(field string, value string) Filter {
	b := wildCard{name: field}
	b.Wildcard = value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) WildcardSuffix(field string, value string) Filter {
	b := wildCard{name: field}
	b.Wildcard = value + "*"
	f.condition = append(f.condition, b)
	return f
}
func (f filter) Result() []interface{} {
	return f.condition
}

func AggDataHistogram(field, interval, format, offset, timeZone string) Agg {
	return agg{}.Name(field).DateHistogram(field, interval, format, offset, timeZone)
}
func AggDataHistogramName(name, field, interval, format, offset, timeZone string) Agg {
	return agg{}.Name(name).DateHistogram(field, interval, format, offset, timeZone)
}

func AggFilters(name string, subAgg Agg, filters ...AggFilter) Agg {
	return agg{}.Name(name).Filter(subAgg, filters...)
}

func AggDistinct(field string, number int64) Agg {
	return agg{}.Name(field).Distinct(field, number)
}

// agg 必须有MarshalJSON，用来生成查询es 需要条件
// 同一个agg中filter 与其他互斥
type agg struct {
	name string
	// 在项目实际使用filter 有多个metric 与 data_histogram 同时使用出现，稳定复现buckets中key 并非预期的问题
	filters       *aggFilters
	dataHistogram *aggDateHistogram `json:"date_histogram"`
	distinct      *aggDistinct      `json:"terms"`
}

func (a agg) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{}, 0)
	if a.filters == nil {
		if a.dataHistogram != nil {
			result["date_histogram"] = a.dataHistogram
		} else if a.distinct != nil {
			result["terms"] = a.distinct
		}
		return json.Marshal(result)
	}

	filters := make(map[string]interface{}, len(a.filters.buckets))
	for _, item := range a.filters.buckets {
		name, cond := item.Result()
		filters[name] = cond
	}
	result["filters"] = map[string]interface{}{"filters": filters}
	subAggName, subAgg := a.filters.agg.Result()
	result["aggs"] = map[string]interface{}{subAggName: subAgg}
	return json.Marshal(result)
}

func (a agg) Name(name string) Agg {
	a.name = name
	return a
}

func (a agg) Filter(agg Agg, filter ...AggFilter) Agg {
	a.filters = &aggFilters{buckets: filter, agg: agg}
	return a
}

func (a agg) DateHistogram(field, interval, format, offset, timeZone string) Agg {
	a.dataHistogram = &aggDateHistogram{
		Field:            field,
		CalendarInterval: interval,
		Format:           format,
		Offset:           offset,
		TimeZone:         timeZone,
	}
	return a
}

func (a agg) Distinct(field string, number int64) Agg {
	a.distinct = &aggDistinct{
		Field: field,
		Size:  number,
	}

	return a
}

func (a agg) Result() (string, interface{}) {
	return a.name, a
}

type aggDateHistogram struct {
	Field string `json:"field,omitempty"`
	// minute,hour,day,week,month,quarter,year
	CalendarInterval string `json:"calendar_interval,omitempty"`
	Format           string `json:"format,omitempty"`
	Offset           string `json:"offset,omitempty"`
	TimeZone         string `json:"time_zone,omitempty"`
}

type aggDistinct struct {
	Field string `json:"field"`
	Size  int64  `json:"size"`
}

//  aggFilter
//  @Description:
//
type aggFilters struct {
	buckets []AggFilter
	agg     Agg
}

func AggFilterName(name string) AggFilter {
	return &aggFilter{
		name:  name,
		terms: nil,
	}
}

type aggFilter struct {
	name  string
	terms map[string]interface{}
}

func (a aggFilter) Name(name string) AggFilter {
	a.name = name
	return a
}

// Terms  raw value must be slice or array
func (a aggFilter) Terms(field string, val ...interface{}) AggFilter {
	if a.terms == nil {
		a.terms = make(map[string]interface{}, 0)
	}
	a.terms[field] = val
	return a
}

func (a aggFilter) TermsArray(field string, arr interface{}) AggFilter {
	if a.terms == nil {
		a.terms = make(map[string]interface{}, 0)
	}
	a.terms[field] = arr
	return a
}

func (a aggFilter) Result() (string, map[string]interface{}) {
	result := make(map[string]interface{}, 0)
	if a.terms != nil {
		result["terms"] = a.terms
	}
	return a.name, result
}

type multiUpdate struct {
	MFilter map[string]interface{}
	MDoc    map[string]interface{}
}

func (m multiUpdate) Filter(filter map[string]interface{}) MultiUpdate {
	if m.MFilter == nil {
		m.MFilter = make(map[string]interface{}, 0)
	}
	m.MFilter = filter
	return m
}

func (m multiUpdate) Doc(doc map[string]interface{}) MultiUpdate {
	if m.MDoc == nil {
		m.MDoc = make(map[string]interface{}, 0)
	}
	m.MDoc = doc
	return m
}

func (m multiUpdate) Get() (map[string]interface{}, map[string]interface{}) {
	return m.MFilter, m.MDoc
}

func MultiUpdateFilter(filter, doc map[string]interface{}) MultiUpdate {
	return multiUpdate{}.Filter(filter).Doc(doc)
}

var (
	_ Agg         = (*agg)(nil)
	_ Filter      = (*filter)(nil)
	_ AggFilter   = (*aggFilter)(nil)
	_ MultiUpdate = (*multiUpdate)(nil)
)
