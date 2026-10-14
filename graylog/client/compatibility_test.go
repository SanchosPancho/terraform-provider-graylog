package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/client"
	"github.com/SanchosPancho/terraform-provider-graylog/graylog/client/system/indices/indexset"
	"github.com/SanchosPancho/terraform-provider-graylog/graylog/config"
)

func TestGraylog7Requests(t *testing.T) {
	ctx := context.Background()
	stream := map[string]interface{}{
		"title": "test", "description": "test", "index_set_id": "index-id",
		"matching_type": "AND", "remove_matches_from_default_stream": false,
	}
	tests := []struct {
		name, version, method, path string
		input, want                 map[string]interface{}
		call                        func(client.Client, map[string]interface{}) error
	}{
		{
			name: "stream creation", version: "v7", method: "POST", path: "/streams",
			input: stream, want: map[string]interface{}{"entity": stream},
			call: func(c client.Client, d map[string]interface{}) error {
				_, _, err := c.Stream.Create(ctx, d)
				return err
			},
		},
		{
			name: "legacy stream creation", version: "v4", method: "POST", path: "/streams",
			input: stream, want: stream,
			call: func(c client.Client, d map[string]interface{}) error {
				_, _, err := c.Stream.Create(ctx, d)
				return err
			},
		},
		{
			name: "strict stream creation", version: "v7", method: "POST", path: "/streams",
			input: map[string]interface{}{"title": "test", "id": "response-id", "disabled": true, "created_at": "date", "creator_user_id": "admin", "is_default": false},
			want:  map[string]interface{}{"entity": map[string]interface{}{"title": "test"}},
			call: func(c client.Client, d map[string]interface{}) error {
				_, _, err := c.Stream.Create(ctx, d)
				return err
			},
		},
		{
			name: "stream update stays flat", version: "v7", method: "PUT", path: "/streams/stream-id",
			input: map[string]interface{}{"title": "updated", "disabled": true, "created_at": "date"},
			want:  map[string]interface{}{"title": "updated"},
			call: func(c client.Client, d map[string]interface{}) error {
				_, _, err := c.Stream.Update(ctx, "stream-id", d)
				return err
			},
		},
		{
			name: "index update stays flat", version: "v7", method: "PUT", path: "/system/indices/index_sets/index-id",
			input: map[string]interface{}{"title": "updated", "rotation_strategy": map[string]interface{}{"type": "rotation"}, "retention_strategy": map[string]interface{}{"type": "retention"}, "data_tiering": nil, "use_legacy_rotation": true, "creation_date": "date", "index_prefix": "prefix", "index_analyzer": "standard", "default": false},
			want:  map[string]interface{}{"id": "index-id", "title": "updated", "rotation_strategy": map[string]interface{}{"type": "rotation"}, "retention_strategy": map[string]interface{}{"type": "retention"}, "data_tiering": nil, "use_legacy_rotation": true},
			call: func(c client.Client, d map[string]interface{}) error {
				_, _, err := c.IndexSet.Update(ctx, "index-id", d)
				return err
			},
		},
		{
			name: "notification creation", version: "v7", method: "POST", path: "/events/notifications",
			input: map[string]interface{}{"title": "test"}, want: map[string]interface{}{"entity": map[string]interface{}{"title": "test"}},
			call: func(c client.Client, d map[string]interface{}) error {
				_, _, err := c.EventNotification.Create(ctx, d)
				return err
			},
		},
		{
			name: "event definition creation", version: "v7", method: "POST", path: "/events/definitions",
			input: map[string]interface{}{"title": "test"}, want: map[string]interface{}{"entity": map[string]interface{}{"title": "test"}},
			call: func(c client.Client, d map[string]interface{}) error {
				_, _, err := c.EventDefinition.Create(ctx, d)
				return err
			},
		},
		{
			name: "view creation", version: "v7", method: "POST", path: "/views",
			input: map[string]interface{}{"title": "test"}, want: map[string]interface{}{"entity": map[string]interface{}{"title": "test"}},
			call: func(c client.Client, d map[string]interface{}) error { _, _, err := c.View.Create(ctx, d); return err },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.method, r.Method)
				assert.Equal(t, "/api"+tt.path, r.URL.Path)
				assert.Equal(t, "terraform-provider-graylog", r.Header.Get("X-Requested-By"))
				name, password, ok := r.BasicAuth()
				assert.True(t, ok)
				assert.Equal(t, "test-token", name)
				assert.Equal(t, "token", password)
				var body map[string]interface{}
				assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				assert.Equal(t, tt.want, body)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{}`))
			}))
			defer server.Close()
			cl, err := client.New(config.Config{Endpoint: server.URL + "/api", AuthName: "test-token", AuthPassword: "token", APIVersion: tt.version})
			require.NoError(t, err)
			require.NoError(t, tt.call(cl, tt.input))
		})
	}
}

func TestGraylog7IndexCreation(t *testing.T) {
	input := map[string]interface{}{"title": "test", "default": false, "id": "response-id", "data_tiering": nil, "use_legacy_rotation": true}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NotContains(t, body, "entity")
		assert.NotContains(t, body, "default")
		assert.NotContains(t, body, "id")
		assert.Equal(t, true, body["use_legacy_rotation"])
		assert.Contains(t, body, "data_tiering")
		assert.Nil(t, body["data_tiering"])
		date, ok := body["creation_date"].(string)
		assert.True(t, ok)
		_, err := time.Parse(time.RFC3339Nano, date)
		assert.NoError(t, err)
		_, _ = w.Write([]byte(`{"id":"index-id"}`))
	}))
	defer server.Close()
	cl, err := client.New(config.Config{Endpoint: server.URL + "/api", APIVersion: "v7"})
	require.NoError(t, err)
	_, _, err = cl.IndexSet.Create(context.Background(), input)
	require.NoError(t, err)
	require.NotContains(t, input, "creation_date", "the caller's map must not be mutated")
}

func TestIndexSetListParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "1", r.URL.Query().Get("skip"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "true", r.URL.Query().Get("stats"))
		_, _ = w.Write([]byte(`{"index_sets":[]}`))
	}))
	defer server.Close()
	cl, err := client.New(config.Config{Endpoint: server.URL + "/api", APIVersion: "v7"})
	require.NoError(t, err)
	_, _, err = cl.IndexSet.Gets(context.Background(), &indexset.GetAllParams{Skip: 1, Limit: 10, Stats: true})
	require.NoError(t, err)
}

func TestGraylog7UserReadUsesID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/users/id/user-id", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"user-id"}`))
	}))
	defer server.Close()
	cl, err := client.New(config.Config{Endpoint: server.URL + "/api", APIVersion: "v7"})
	require.NoError(t, err)
	_, _, err = cl.User.Get(context.Background(), "user-id", cl.APIVersion)
	require.NoError(t, err)
}
