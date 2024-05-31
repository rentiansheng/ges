package ges

/***************************
    @author: tiansheng.ren
    @date: 2024/5/29
    @desc:

***************************/

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
	_ NestedFilter = (*esNested)(nil)
)
