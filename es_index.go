package ges

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/12/16
    @desc:
		func name with prefix I,  Informal API. eg: ICreate,

***************************/

type esIndex struct {
	name string
}

func (e esIndex) Exists(ctx context.Context) (bool, error) {
	res, err := rawESClient.Indices.Exists([]string{e.name})
	if err != nil {
		return false, err
	}
	if res.StatusCode == 404 {
		return false, nil
	}
	if res.StatusCode == 200 {
		return true, nil
	}
	if _, err := parseSearchRespDefaultDecode(ctx, res); err != nil {
		return false, err
	}

	return false, nil
}

func (e esIndex) Create(ctx context.Context, mapping IndexMeta) error {

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(mapping); err != nil {
		return fmt.Errorf("index create mapping encode error. %s", err.Error())
	}
	res, err := rawESClient.Indices.Create(e.name, esapi.IndicesCreate.WithBody(nil, body))
	if err != nil {
		return fmt.Errorf("es client do error. %s", err.Error())
	}

	if _, err := parseSearchRespDefaultDecode(ctx, res); err != nil {
		return err
	}
	return nil
}

func (e esIndex) List(ctx context.Context) ([]error, error) {
	//TODO implement me
	panic("implement me")
}

func (e esIndex) Mapping(ctx context.Context) (map[string]indexMetaResp, error) {
	res, err := rawESClient.Indices.Get([]string{e.name})
	if err != nil {
		return nil, err
	}
	indexRes, err := parseSearchRespIndexDecode(ctx, res)
	if err != nil {
		return nil, err
	}

	return indexRes, nil
}

var _ Index = (*esIndex)(nil)
