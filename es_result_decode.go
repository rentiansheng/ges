package ges

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/12/16
    @desc:

***************************/

func parseSearchRespDefaultDecode(ctx context.Context, res *esapi.Response) (SearchResult, error) {
	var resp SearchResult

	if res.IsError() {
		return resp, fmt.Errorf("elasticsearch response error. status: %v, message: %s", res.Status(), res.String())
	}

	respBody := res.Body

	d := json.NewDecoder(respBody)
	d.UseNumber()
	err := d.Decode(&resp)
	if err != nil {
		return resp, err
	}
	if resp.Error != nil {
		return resp, fmt.Errorf("%s", resp.Error)
	}
	if resp.TimeOut {
		return resp, fmt.Errorf(" time_out, took: %v", resp.Took)
	}
	return resp, nil
}

func parseSearchRespIndexDecode(ctx context.Context, res *esapi.Response) (map[string]indexMetaResp, error) {
	resp := make(map[string]indexMetaResp, 1)

	if res.IsError() {
		return resp, fmt.Errorf("elasticsearch response error. status: %v, message: %s", res.Status(), res.String())
	}
	if res.StatusCode != 200 {
		return resp, fmt.Errorf("code: %v, message: %s", res.StatusCode, res.Status())
	}
	respBody := res.Body

	d := json.NewDecoder(respBody)
	d.UseNumber()
	err := d.Decode(&resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
