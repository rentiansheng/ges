package ges

import (
	"context"
	"encoding/json"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/7/27
    @desc:

***************************/

const (
	MaxBulkItemsLimit       = 1000
	BulkItemsLimit          = 100
	MaxBulkUpdateItemsLimit = 10
)

type Client interface {
	IndexName(name string) Client
	Index() Index

	Not(filters ...Filter) Client
	Where(filters ...Filter) Client
	Or(filters ...Filter) Client
	OrderBy(field string, isDesc bool) Client
	Size(uint64) Client
	Agg(aggs ...Agg) Client
	Start(uint64) Client
	Limit(uint64, uint64) Client
	Limit64(int64, int64) Client
	Fields(...string) Client
	Search(ctx context.Context, result interface{}) (uint64, error)
	GetById(ctx context.Context, id string, result interface{}) error
	RawSQL(ctx context.Context, sql string, result interface{}) error
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
	TranslateSQL(ctx context.Context, sql string) ([]byte, error)
	Query(ctx context.Context, raw interface{}, result interface{}) error
}

type Filter interface {
	Term(field string, value interface{}) Filter
	Terms(field string, values interface{}) Filter
	TermsSingeItem(field string, value interface{}) Filter
	Between(field string, start, end int64) Filter
	Gt(field string, value int64) Filter
	Gte(field string, value int64) Filter
	Lt(field string, value int64) Filter
	Lte(field string, value int64) Filter
	Wildcard(field string, value string) Filter
	WildcardSuffix(field string, value string) Filter
	Result() []interface{}
}

type Index interface {
	Exists(ctx context.Context) (bool, error)
	Create(ctx context.Context, mapping IndexMeta) error
	List(ctx context.Context) ([]error, error)
	Mapping(ctx context.Context) (map[string]indexMetaResp, error)
}

type Agg interface {
	Name(string) Agg
	DateHistogram(field, interval, format, offset, timeZone string) Agg
	DateHistogramAgg(field, interval, format, offset, timeZone string, sub ...Agg) Agg
	Distinct(field string, number int64) Agg
	Filter(agg Agg, filter ...AggFilter) Agg
	Sum(field string) Agg
	Nested(path string) Agg
	Avg(field string) Agg
	//Metric(...AggDateHistogramMetric) Agg
	Result() (string, interface{})
}

type AggFilter interface {
	Name(string) AggFilter
	Terms(field string, val ...interface{}) AggFilter
	TermsArray(field string, val interface{}) AggFilter

	Result() (string, map[string]interface{})
}

type AggDateHistogramMetric interface {
	Add(name, operator string)
	Result() interface{}
}

// SearchResult index 返回数据，直接解析到对应的结构体
type SearchResult struct {
	Took    uint64       `json:"took"`
	TimeOut bool         `json:"time_out"`
	Error   interface{}  `json:"error"`
	Shards  ShardsResult `json:"_shards"`
	Hits    struct {
		Total struct {
			Value    uint64 `json:"value"`
			Relation string `json:"relation"`
		}
		MaxScore  float64                 `json:"max_score"`
		IndexHits []SearchResultHitResult `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations"`
}

// SearchResultHitResult index 返回数据，直接解析到对应的结构体
type SearchResultHitResult struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	Id     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

type CountResult struct {
	Error  interface{}  `json:"error"`
	Count  uint64       `json:"count"`
	Shards ShardsResult `json:"_shards"`
}

type ShardsResult struct {
	Total      uint64 `json:"total"`
	Successful uint64 `json:"successful"`
	Skipped    uint64 `json:"skipped"`
	Failed     uint64 `json:"failed"`
}

type SQLResult struct {
	Error   interface{}     `json:"error"`
	Columns []SQLColumn     `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type SQLCountResult struct {
	Error   interface{} `json:"error"`
	Columns []SQLColumn `json:"columns"`
	Rows    [][]uint64  `json:"rows"`
}

func (s SQLCountResult) Count() uint64 {
	if len(s.Rows) > 0 && len(s.Rows[0]) > 0 {
		return s.Rows[0][0]
	}
	return 0
}

type SQLColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type MultiUpdate interface {
	Filter(filter map[string]interface{}) MultiUpdate
	Doc(doc map[string]interface{}) MultiUpdate
	Get() (filter map[string]interface{}, doc map[string]interface{})
}

type IndexMeta struct {
	Settings *IndexMappingSettings `json:"settings,omitempty"`
	Mappings IndexMapping          `json:"mappings"`
}

type IndexMapping struct {
	Properties map[string]MappingField `json:"properties"`
}

type MappingField struct {
	Type       MappingType             `json:"type"`
	Properties map[string]MappingField `json:"properties,omitempty"`
}

type MappingType string

const (
	MappingTypeBinary  MappingType = "binary"
	MappingTypeBoolean MappingType = "boolean"
	MappingTypeKeyword MappingType = "keyword"
	// long A signed 64-bit integer with a minimum value of -263 and a maximum value of 263-1.
	MappingTypeLong MappingType = "long"
	// MappingTypeInteger A signed 32-bit integer with a minimum value of -231 and a maximum value of 231-1
	MappingTypeInteger MappingType = "integer"
	// MappingTypeShort A signed 16-bit integer with a minimum value of -32,768 and a maximum value of 32,767
	MappingTypeShort MappingType = "short"
	// MappingTypeByte A signed 8-bit integer with a minimum value of -128 and a maximum value of 127.
	MappingTypeByte MappingType = "byte"
	// MappingTypeDouble A double-precision 64-bit IEEE 754 floating point number, restricted to finite values.
	MappingTypeDouble MappingType = "double"
	// MappingTypeFloat A single-precision 32-bit IEEE 754 floating point number, restricted to finite values.
	MappingTypeFloat MappingType = "float"
	// MappingTypeHalfFloat  half-precision 16-bit IEEE 754 floating point number, restricted to finite values.
	MappingTypeHalfFloat MappingType = "half_float"
	// MappingTypeScaledFloat A floating point number that is backed by a long, scaled by a fixed double scaling factor
	MappingTypeScaledFloat MappingType = "scaled_float"
	// MappingTypeUnsignedLong An unsigned 64-bit integer with a minimum value of 0 and a maximum value of 264-1
	MappingTypeUnsignedLong MappingType = "unsigned_long"
	MappingTypeDate         MappingType = "date"
	MappingTypeNested       MappingType = "nested"
	MappingTypeText         MappingType = "text"
)

type indexMetaResp struct {
	Aliases  map[string]interface{} `json:"aliases"`
	Settings IndexMappingSettings   `json:"settings"`
	Mappings IndexMetaRespMappings  `json:"mappings"`
}

type IndexMappingSettings struct {
	Index IndexSettings `json:"index,omitempty"`
}

type IndexSettings struct {
	CreationDate     string `json:"creation_date,omitempty"`
	NumberOfShards   int    `json:"number_of_shards,omitempty"`
	NumberOfReplicas int    `json:"number_of_replicas,omitempty"`
	Uuid             string `json:"uuid,omitempty"`
	Version          *struct {
		Created string `json:"created,omitempty"`
	} `json:"version,omitempty"`
	ProvidedName string `json:"provided_name,omitempty"`
}

type IndexMetaRespMappings struct {
	Properties map[string]IndexMetaRespProperties `json:"properties"`
}

type IndexMetaRespProperties struct {
	Properties IndexMetaRespPropertiesFieldProperties `json:"properties"`
}

type IndexMetaRespPropertiesFieldProperties struct {
	Type IndexMetaRespPropertiesFieldPropertiesField `json:"type"`
}

type IndexMetaRespPropertiesFieldPropertiesField struct {
	Type MappingType `json:"type"`
	// 不同数据类型题不同
	Fields map[string]interface{} `json:"fields,omitempty"`
}

type SourceResp struct {
	Index       string          `json:"_index"`
	Type        string          `json:"_type"`
	Id          string          `json:"_id"`
	Version     int             `json:"_version"`
	SeqNo       int             `json:"_seq_no"`
	PrimaryTerm int             `json:"_primary_term"`
	Found       bool            `json:"found"`
	Source      json.RawMessage `json:"_source"`
	Error       interface{}     `json:"error"`
	Status      int             `json:"status"`
}

type Document interface {
	ID() string
	Doc() interface{}
	Item() (string, interface{})
}
