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

var (
	_ Filter    = (*filter)(nil)
	_ AggFilter = (*aggFilter)(nil)
)
