package ges

import (
	"encoding/json"
)

/***************************
    @author: tiansheng.ren
    @date: 2024/5/29
    @desc:

***************************/

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
