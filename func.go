package ges

/***************************
    @author: tiansheng.ren
    @date: 2024/5/29
    @desc:

***************************/

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

func Append(filters ...Filter) Filter {
	return filter{}.Append(filters...)
}
