package indexset

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/client/api"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client     httpclient.Client
	APIVersion string
}

func (cl Client) Get(ctx context.Context, id string) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/indices/index_sets/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

type GetAllParams struct {
	Skip  int
	Limit int
	Stats bool
}

func (cl Client) Gets(
	ctx context.Context, params *GetAllParams,
) (map[string]interface{}, *http.Response, error) {
	query := url.Values{}
	if params != nil {
		if params.Skip != 0 {
			query.Add("skip", strconv.Itoa(params.Skip))
		}
		if params.Limit != 0 {
			query.Add("limit", strconv.Itoa(params.Limit))
		}
		if params.Stats {
			query.Add("stats", "true")
		}
	}
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/indices/index_sets",
		Query:        query,
		ResponseBody: &body,
	})
	return body, resp, err
}

func genCreationDate() string {
	return time.Now().In(time.FixedZone("UTC", 0)).Format(time.RFC3339Nano)
}

func (cl Client) Create(
	ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}
	if cl.APIVersion == "v7" {
		data = api.Fields(data, "title", "description", "index_prefix", "index_analyzer", "shards", "replicas", "writable", "rotation_strategy_class", "rotation_strategy", "retention_strategy_class", "retention_strategy", "index_optimization_max_num_segments", "index_optimization_disabled", "field_type_refresh_interval", "use_legacy_rotation", "data_tiering", "field_type_profile", "index_template_type", "index_set_template_id", "field_restrictions", "creation_date")
	}

	body := map[string]interface{}{}
	if v, ok := data["creation_date"]; !ok || v == "" || v == nil {
		data["creation_date"] = genCreationDate()
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/system/indices/index_sets",
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Update(
	ctx context.Context, id string, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}
	if cl.APIVersion == "v7" {
		data = api.Fields(data, "title", "description", "shards", "replicas", "writable", "rotation_strategy_class", "rotation_strategy", "retention_strategy_class", "retention_strategy", "index_optimization_max_num_segments", "index_optimization_disabled", "field_type_refresh_interval", "use_legacy_rotation", "data_tiering", "field_type_profile", "field_restrictions")
		data["id"] = id
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "PUT",
		Path:         "/system/indices/index_sets/" + id,
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Delete(ctx context.Context, id string) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/system/indices/index_sets/" + id,
	})
	return resp, err
}
