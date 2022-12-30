package ges

/***************************
    @author: tiansheng.ren
    @date: 2022/11/5
    @desc:

***************************/

type DateHistogramBuckets struct {
	Buckets []DateHistogramBucket `json:"buckets"`
}

type DateHistogramBucket struct {
	Key    int64   `json:"key"`
	StrKey string  `json:"str_key"`
	Count  float64 `json:"doc_count"`
}

type DistinctBuckets struct {
	Buckets []DistinctBucket `json:"buckets"`
}

type DistinctBucket struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}

const (
	/** date histogram interval type  https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-datehistogram-aggregation.html **/

	DateHistogramIntervalMinute  = "minute"
	DateHistogramIntervalHour    = "hour"
	DateHistogramIntervalDay     = "day"
	DateHistogramIntervalWeek    = "week"
	DateHistogramIntervalMonth   = "month"
	DateHistogramIntervalQuarter = "quarter"
	DateHistogramIntervalYear    = "year"
)
