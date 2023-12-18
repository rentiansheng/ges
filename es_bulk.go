package ges

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/11/28
    @desc:

***************************/

type bulkItemDetailResp struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	Id     string          `json:"_id"`
	Result string          `json:"result"`
	Error  json.RawMessage `json:"error"`
}

type bulkItemResp struct {
	Index  bulkItemDetailResp `json:"index"`
	Update bulkItemDetailResp `json:"update"`
	Create bulkItemDetailResp `json:"create"`
	DELETE bulkItemDetailResp `json:"delete"`
}

type bulkResp struct {
	Took   int            `json:"took"`
	Errors bool           `json:"errors"`
	Items  []bulkItemResp `json:"items"`
}

func parseBulkResp(ctx context.Context, res *esapi.Response) (resp *bulkResp, err error) {

	respBody := res.Body

	// http status_code not 2xx
	if res.IsError() {
		return nil, fmt.Errorf("bulk fail, status_code: %d, body: %s", res.StatusCode, res.String())
	}

	var r bulkResp
	if err := json.NewDecoder(respBody).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s, %s", err, res.String())
	}
	// bulk error
	if r.Errors {
		jsonR, _ := json.Marshal(r.Items)
		return nil, fmt.Errorf("bulk fail: %+v, %+v, %+v", r.Took, r.Errors, string(jsonR))
	}
	return &r, nil
}
