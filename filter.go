package ges

import (
	"encoding/json"
)

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
	values interface{}
}

func (t terms) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]interface{}{
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
	LtePtr interface{} `json:"lte,omitempty"`
	LtPtr  interface{} `json:"lt,omitempty"`
	GtePtr interface{} `json:"gte,omitempty"`
	GtPtr  interface{} `json:"gt,omitempty"`
}

// FromTo  range query. [from, to), from <= field < to
func (f betweenValue) FromTo(from, to interface{}) RangeFilter {
	f.GtePtr, f.LtPtr = from, to
	return f
}

func (f betweenValue) Range(start, end interface{}) RangeFilter {
	f.GtePtr, f.LtePtr = start, end
	return f
}

func (f betweenValue) Between(start, end interface{}) RangeFilter {
	f.GtePtr, f.LtePtr = start, end
	return f
}

func (f betweenValue) Gt(value interface{}) RangeFilter {
	f.GtPtr = value
	return f
}

func (f betweenValue) Gte(value interface{}) RangeFilter {
	f.GtePtr = value
	return f
}

func (f betweenValue) Lt(value interface{}) RangeFilter {
	f.LtPtr = value
	return f
}

func (f betweenValue) Lte(value interface{}) RangeFilter {

	f.LtePtr = value
	return f
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

type boolFilter struct {
	Bool esQueryBool `json:"bool"`
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
	Must               []interface{} `json:"must,omitempty"`
	Not                []interface{} `json:"must_not,omitempty"`
	Should             []interface{} `json:"should,omitempty"`
	Match              []interface{} `json:"match,omitempty"`
	AdjustPureNegative bool          `json:"adjust_pure_negative,omitempty"`
}

type filter struct {
	condition []interface{}
}

func NewFilter() Filter {
	return filter{}
}

func Term(field string, value interface{}) Filter { return filter{}.Term(field, value) }
func Terms(field string, values interface{}) Filter {
	return filter{}.Terms(field, values)
}
func TermsSingeItem(field string, value interface{}) Filter {
	return filter{}.TermsSingeItem(field, value)
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
func Bool(must, not, should, match Filter, adjustPureNegative bool) Filter {
	return filter{}.BoolItem(must, not, should, match, adjustPureNegative)
}

// Range  range query. [start, end], start <= field <= end
func Range(field string, start, end interface{}) Filter { return filter{}.Range(field, start, end) }

// FromTo  range query. [from, to), from <= field < to
func FromTo(field string, from, to interface{}) Filter { return filter{}.FromTo(field, from, to) }

func NestedQuery(path string, must, not, should, match Filter) NestedFilter {
	return esNested{}.Path(path).Must(must).Not(not).Should(should).Match(match)
}

func NewRangeFilter() RangeFilter {
	return betweenValue{}
}

func Nested() NestedFilter {
	return esNested{}
}

func BoolTrue(must, not, should, match Filter) Filter {
	return filter{}.BoolItem(must, not, should, match, true)
}

func (f filter) Term(field string, value interface{}) Filter {
	f.condition = append(f.condition, term{field, value})
	return f
}

func (f filter) TermsSingeItem(field string, value interface{}) Filter {
	f.condition = append(f.condition, terms{field, []interface{}{value}})
	return f
}

func (f filter) Terms(field string, values interface{}) Filter {
	f.condition = append(f.condition, terms{field, values})
	return f
}

// FromTo  range query. [from, to), from <= field < to
func (f filter) FromTo(field string, from, to interface{}) Filter {
	b := between{name: field}

	b.GtePtr, b.LtPtr = from, to
	f.condition = append(f.condition, b)
	return f
}

// Range  range query. [start, end], start <= field <= end
func (f filter) Range(field string, start, end interface{}) Filter {
	b := between{name: field}
	b.GtePtr, b.LtePtr = start, end
	f.condition = append(f.condition, b)
	return f
}

// Between  range query. [start, end], start <= field <= end
func (f filter) Between(field string, start, end int64) Filter {
	b := between{name: field}
	b.GtePtr, b.LtePtr = start, end
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Gt(field string, value int64) Filter {
	b := between{name: field}
	b.GtPtr = value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Gte(field string, value int64) Filter {
	b := between{name: field}
	b.GtePtr = value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Lt(field string, value int64) Filter {
	b := between{name: field}
	b.LtPtr = value
	f.condition = append(f.condition, b)
	return f
}

func (f filter) Lte(field string, value int64) Filter {
	b := between{name: field}

	b.LtePtr = value
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

// BoolItem 用于构建bool查询, 后期需要优化，暴露出来 bool query 对象，用来管理条件，
func (f filter) BoolItem(must, not, should, match Filter, adjustPureNegative bool) Filter {
	if must == nil && not == nil && should == nil && match == nil {
		return f
	}

	b := esQueryBool{

		AdjustPureNegative: adjustPureNegative,
	}
	if must != nil {
		b.Must = must.Result()
	}
	if not != nil {
		b.Not = not.Result()
	}
	if should != nil {
		b.Should = should.Result()
	}
	if match != nil {
		b.Match = match.Result()
	}

	f.condition = append(f.condition, boolFilter{Bool: b})
	return f
}

func (f filter) Nested(nestedFilter NestedFilter) Filter {
	f.condition = append(f.condition, nestedFilter.Result())
	return f
}

func (f filter) Result() []interface{} {
	return f.condition
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

type esNested struct {
	NestedPath string        `json:"path"`
	Query      esNestedQuery `json:"query"`
}

type esNestedQuery struct {
	Bool esQueryBool `json:"bool"`
}

func (e esNested) Must(filters ...Filter) NestedFilter {
	for _, filter := range filters {
		e.Query.Bool.Must = append(e.Query.Bool.Must, filter.Result()...)
	}
	return e
}

func (e esNested) Should(filters ...Filter) NestedFilter {
	for _, filter := range filters {
		e.Query.Bool.Should = append(e.Query.Bool.Should, filter.Result()...)
	}
	return e
}

func (e esNested) Not(filters ...Filter) NestedFilter {
	for _, filter := range filters {
		e.Query.Bool.Not = append(e.Query.Bool.Not, filter.Result()...)
	}
	return e
}

func (e esNested) Match(filters ...Filter) NestedFilter {
	for _, filter := range filters {
		e.Query.Bool.Match = append(e.Query.Bool.Match, filter.Result()...)
	}
	return e
}

func (e esNested) Path(path string) NestedFilter {
	e.NestedPath = path
	return e
}

func (e esNested) Result() interface{} {
	return map[string]interface{}{"nested": e}
}

var (
	_ Agg          = (*agg)(nil)
	_ Filter       = (*filter)(nil)
	_ AggFilter    = (*aggFilter)(nil)
	_ MultiUpdate  = (*multiUpdate)(nil)
	_ NestedFilter = (*esNested)(nil)
)
