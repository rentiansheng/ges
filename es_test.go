package ges

import (
	"context"
	"fmt"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"

	esMock "github.com/orlangure/gnomock/preset/elastic"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/12/16
    @desc:

***************************/

var (
	ctx = context.TODO()
)

func esInit(t *testing.T) (func(), error) {
	fn := func() {}
	//
	tmpInit := esMock.Preset(esMock.WithVersion("7.8.1"))
	esContainer, err := gnomock.Start(tmpInit)
	if err != nil {
		return nil, err
	}
	addr := "http://" + esContainer.DefaultAddress()
	fmt.Println("es addr:", addr)
	fn = func() {
		err := gnomock.Stop(esContainer)
		require.NoError(t, err, "stop elasticsearch")
	}
	if err := InitClientWithCfg([]string{addr}, "", ""); err != nil {
		return fn, err
	}
	return fn, nil
}

func esInitTest(t *testing.T) (func(), error) {
	deferFn, err := esInit(t)
	if err != nil {
		return deferFn, err
	}
	es := ES().IndexName(indexName)
	esIdx := es.Index()
	err = esIdx.Create(ctx, indexMeta)
	if err != nil {
		return deferFn, err
	}
	return deferFn, nil
}

type testIndexMappingRow struct {
	EsId  string `json:"_id"`
	Id    int64  `json:"tid"`
	Label string `json:"label"`
}

var (
	indexName = "test_index"
	indexMeta = IndexMeta{
		Mappings: IndexMapping{
			Properties: map[string]MappingField{
				"tid":   {Type: MappingTypeInteger},
				"label": {Type: MappingTypeKeyword},
			},
		},
	}
)

func TestESIndexCreateAndCountAndUpdate(t *testing.T) {
	deferFn, err := esInit(t)
	defer deferFn()
	require.NoError(t, err, "es init error")

	es := ES().IndexName(indexName)
	esIdx := es.Index()
	err = esIdx.Create(ctx, indexMeta)
	require.NoError(t, err, "es create index error")

	docs := []interface{}{
		mapStrAny{"tid": int32(1), "label": "tid-1"},
		mapStrAny{"tid": 2, "label": "tid-2"},
		mapStrAny{"tid": 3, "label": "tid-3"},
	}
	err = es.Save(ctx, docs...)
	require.NoError(t, err, "es save docs error")

	cnt, err := es.Count(ctx)
	require.NoError(t, err, "es count docs error ")
	require.Equal(t, len(docs), int(cnt), "es count docs error")

	rows := make([]testIndexMappingRow, 0, 3)
	cnt, err = es.Where(Terms("tid", 1)).Search(ctx, &rows)
	require.NoError(t, err, "es search docs error")
	require.Equal(t, 1, len(rows), "es search docs error")
	if cnt == 0 {
		return
	}

	esId := rows[0].EsId
	newLabelVal := "tid-1-update"
	err = es.UpdateById(ctx, esId, map[string]interface{}{
		"label": newLabelVal,
	})
	require.NoError(t, err, "es search docs error")

	esRow := testIndexMappingRow{}
	err = es.GetById(ctx, esId, &esRow)
	require.NoError(t, err, "es get by id error")
	require.Equal(t, newLabelVal, esRow.Label, "es get by id error")

	err = es.Delete(ctx)
	require.NoError(t, err, "delete by condition")

	cnt, err = es.Count(ctx)
	require.NoError(t, err, "es delete after count docs error")
	require.Equal(t, 0, int(cnt), "es delete after count docs error")

	err = es.GetById(ctx, esId, &esRow)
	require.EqualError(t, err, NotFoundError.Error(), "es get by id error")

}

func TestESUpsertAndMUpsertAndMUpdate(t *testing.T) {
	deferFn, err := esInitTest(t)
	defer deferFn()
	require.NoError(t, err, "es init error")

	doc := mapStrAny{"tid": int64(1), "label": "upsert"}
	esId := "test_upsert_tid_i"

	es := ES().IndexName(indexName)
	err = es.UpsertById(ctx, esId, doc)
	require.NoError(t, err, "es upsert docs error")

	esRow := testIndexMappingRow{}
	err = es.GetById(ctx, esId, &esRow)
	require.NoError(t, err, "es get by id error")
	require.Equal(t, doc["label"], esRow.Label, "es get by id error")

	esId2 := "test_upsert_tid_i"
	docs := mapStrAny{
		esId:  mapStrAny{"tid": int64(1), "label": "multi-upsert-1"},
		esId2: mapStrAny{"tid": int64(2), "label": "multi-upsert-2"},
	}
	err = es.MUpsertById(ctx, DocsFromMap(docs)...)
	require.NoError(t, err, "es multi-upsert docs error")

	rows := make([]testIndexMappingRow, 0)
	cnt, err := es.Search(ctx, &rows)
	require.NoError(t, err, "es search docs error")
	require.Equal(t, len(docs), int(cnt), "es search docs error")

	for _, row := range rows {
		docMap, _ := docs[row.EsId].(mapStrAny)
		require.Equal(t, docMap["tid"], row.Id, "compare multi-upsert result id ")
		require.Equal(t, docMap["label"], row.Label, "compare multi-upsert result label ")
	}

	docs = mapStrAny{
		esId:  mapStrAny{"tid": int64(1), "label": "multi-upset-1"},
		esId2: mapStrAny{"tid": int64(2), "label": "multi-upset-2"},
	}
	err = es.MUpsertById(ctx, DocsFromMap(docs)...)
	require.NoError(t, err, "es multi-upset docs error")

	rows = make([]testIndexMappingRow, 0)
	cnt, err = es.Search(ctx, &rows)
	require.NoError(t, err, "es search docs error")
	require.Equal(t, len(docs), int(cnt), "es search docs error")

	for _, row := range rows {
		docMap, _ := docs[row.EsId].(mapStrAny)
		require.Equal(t, docMap["tid"], row.Id, "compare multi-upsert result id ")
		require.Equal(t, docMap["label"], row.Label, "compare multi-upsert result label ")
	}

}

func TestESUpsertNoID(t *testing.T) {
	deferFn, err := esInitTest(t)
	defer deferFn()
	require.NoError(t, err, "es init error")

	doc := mapStrAny{"tid": int64(1), "label": "upsert not id"}
	esId := ""

	es := ES().IndexName(indexName)
	err = es.UpsertById(ctx, esId, doc)
	require.NoError(t, err, "es upsert docs error")

	rows := make([]testIndexMappingRow, 0)
	cnt, err := es.Search(ctx, &rows)
	require.NoError(t, err, "es search docs error")
	require.Equal(t, 1, int(cnt), "es search docs error")
	for _, row := range rows {
		require.Equal(t, doc["tid"], row.Id, "compare upsert not id result id ")
		require.Equal(t, doc["label"], row.Label, "compare upsert not id result label ")
	}

}

func TestESUSave(t *testing.T) {
	deferFn, err := esInitTest(t)
	defer deferFn()
	require.NoError(t, err, "es init error")

	doc := mapStrAny{"tid": int64(1), "label": "upsert not id"}
	esId := ""

	es := ES().IndexName(indexName)
	err = es.USave(ctx, NewDoc(esId, doc))
	require.NoError(t, err, "es upsert docs error")

	rows := make([]testIndexMappingRow, 0)
	cnt, err := es.Search(ctx, &rows)
	require.NoError(t, err, "es search docs error")
	require.Equal(t, 1, int(cnt), "es search docs error")
	for _, row := range rows {
		require.Equal(t, doc["tid"], row.Id, "compare upsert not id result id ")
		require.Equal(t, doc["label"], row.Label, "compare upsert not id result label ")
	}

}

func TestDelete(t *testing.T) {
	deferFn, err := esInitTest(t)
	defer deferFn()
	require.NoError(t, err, "es init error")

	doc := mapStrAny{"tid": int64(1), "label": "upsert not id"}
	esId := ""

	es := ES().IndexName(indexName)
	err = es.USave(ctx, NewDoc(esId, doc), NewDoc(esId, doc))
	require.NoError(t, err, "es upsert docs error")

	rows := make([]testIndexMappingRow, 0)
	cnt, err := es.Search(ctx, &rows)
	require.NoError(t, err, "es search docs error")
	require.Equal(t, 2, int(cnt), "es search docs error")

	ids := make([]string, 0, 2)
	for _, row := range rows {
		ids = append(ids, row.EsId)
	}
	err = es.DeleteById(ctx, ids...)
	require.NoError(t, err, "es delete by id error")

	cnt, err = es.Count(ctx)
	require.NoError(t, err, "es count docs error")
	require.Equal(t, uint64(0), cnt, "es count docs error")

}
