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
	DocCountErrorUpperBound int              `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int              `json:"sum_other_doc_count"`
	Buckets                 []DistinctBucket `json:"buckets"`
}

type DistinctBucket struct {
	Key      interface{} `json:"key"`
	DocCount int         `json:"doc_count"`
}

type DocCountBuckets struct {
	Buckets map[string]DocCountBucket `json:"buckets"`
}

type DocCountBucket struct {
	DocCount int `json:"doc_count"`
}

type DistinctBucketsI64 struct {
	DocCountErrorUpperBound int                    `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int                    `json:"sum_other_doc_count"`
	Buckets                 []DistinctBucketKeyI64 `json:"buckets"`
}

type DistinctBucketKeyI64 struct {
	Key      int64 `json:"key"`
	DocCount int   `json:"doc_count"`
}

type DistinctBucketsU64 struct {
	DocCountErrorUpperBound int                    `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int                    `json:"sum_other_doc_count"`
	Buckets                 []DistinctBucketKeyU64 `json:"buckets"`
}

type DistinctBucketKeyU64 struct {
	Key      uint64 `json:"key"`
	DocCount int    `json:"doc_count"`
}
