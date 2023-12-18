package ges

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/12/14
    @desc:

***************************/

func TestAggDataHistogram(t *testing.T) {

	agg := AggDataHistogram("field", "week", "yyyy-MM-dd", "+3d", "+08:00")
	name, resultAgg := agg.Result()
	require.Equal(t, "field", name, "TestAggDataHistogram ")
	actual, err := json.Marshal(resultAgg)
	require.NoError(t, err, "TestAggDataHistogram json.Marshal")
	expected := `{"date_histogram":{"field":"field","calendar_interval":"week","format":"yyyy-MM-dd","offset":"+3d","time_zone":"+08:00"}}`
	require.Equal(t, expected, string(actual), "TestAggDataHistogram ")
}

func TestAggDataHistogramName(t *testing.T) {
	agg := AggDataHistogramName("name", "field", "week", "yyyy-MM-dd", "+3d", "+08:00")
	name, resultAgg := agg.Result()
	require.Equal(t, "name", name, "TestAggDataHistogramName ")
	actual, err := json.Marshal(resultAgg)
	require.NoError(t, err, "TestAggDataHistogramName json.Marshal")
	expected := `{"date_histogram":{"field":"field","calendar_interval":"week","format":"yyyy-MM-dd","offset":"+3d","time_zone":"+08:00"}}`
	require.Equal(t, expected, string(actual), "TestAggDataHistogramName ")

}

func TestAggDistinct(t *testing.T) {
	agg := AggDistinct("field", 10)
	name, resultAgg := agg.Result()
	require.Equal(t, "field", name, "TestAggDistinct ")
	actual, err := json.Marshal(resultAgg)
	require.NoError(t, err, "TestAggDistinct json.Marshal")
	expected := `{"terms":{"field":"field","size":10}}`
	require.Equal(t, expected, string(actual), "TestAggDistinct ")

}

func TestAggFilters(t *testing.T) {
	subAgg := AggDataHistogramName("name", "field", "week", "yyyy-MM-dd", "+3d", "+08:00")
	aggFilterArr := []AggFilter{AggFilterName("filter1").Terms("key", "value1", "value2"), AggFilterName("filter2").TermsArray("key", []string{"value"})}
	aggFilters := AggFilters("filter", subAgg, aggFilterArr...)
	name, resultAgg := aggFilters.Result()
	require.Equal(t, "filter", name, "TestAggFilters ")
	actual, err := json.Marshal(resultAgg)
	require.NoError(t, err, "TestAggFilters json.Marshal")
	expected := `{"aggs":{"name":{"date_histogram":{"field":"field","calendar_interval":"week","format":"yyyy-MM-dd","offset":"+3d","time_zone":"+08:00"}}},"filters":{"filters":{"filter1":{"terms":{"key":["value1","value2"]}},"filter2":{"terms":{"key":["value"]}}}}}`
	require.Equal(t, expected, string(actual), "TestAggFilters ")

}

func TestAggDistinctAgg(t *testing.T) {
	agg := AggDistinct("field", 10)
	subAgg := AggDataHistogram("field", "week", "yyyy-MM-dd", "+3d", "+08:00")
	agg = agg.Aggs(subAgg)
	name, resultAgg := agg.Result()
	require.Equal(t, "field", name, "TestAggDistinctAgg ")
	actual, err := json.Marshal(resultAgg)
	require.NoError(t, err, "TestAggDistinctAgg json.Marshal")
	expected := `{"aggs":{"field":{"date_histogram":{"field":"field","calendar_interval":"week","format":"yyyy-MM-dd","offset":"+3d","time_zone":"+08:00"}}},"terms":{"field":"field","size":10}}`
	require.Equal(t, expected, string(actual), "TestAggDistinctAgg ")

}
