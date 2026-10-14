package graylog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/config"
	"github.com/SanchosPancho/terraform-provider-graylog/graylog/resource/stream"
	"github.com/SanchosPancho/terraform-provider-graylog/graylog/resource/stream/rule"
	indexset "github.com/SanchosPancho/terraform-provider-graylog/graylog/resource/system/indices/indexset"
)

func TestGraylog7ResourceLifecycle(t *testing.T) {
	for _, legacy := range []bool{true, false} {
		name := "legacy_rotation"
		if !legacy {
			name = "data_tiering"
		}
		t.Run(name, func(t *testing.T) {
			var indexState, streamState, ruleState map[string]interface{}
			var calls []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := strings.TrimPrefix(r.URL.Path, "/api")
				calls = append(calls, r.Method+" "+path)
				w.Header().Set("Content-Type", "application/json")
				var body map[string]interface{}
				if r.Method == "PUT" || (r.Method == "POST" && !strings.HasSuffix(path, "/resume")) {
					assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				}
				respond := func(value interface{}) { assert.NoError(t, json.NewEncoder(w).Encode(value)) }
				switch r.Method + " " + path {
				case "POST /system/indices/index_sets":
					assert.NotContains(t, body, "entity")
					assert.NotContains(t, body, "default")
					assert.Contains(t, body, "rotation_strategy")
					assert.Contains(t, body, "retention_strategy")
					assert.Equal(t, legacy, body["use_legacy_rotation"])
					if legacy {
						assert.Nil(t, body["data_tiering"])
					} else {
						assert.Equal(t, "hot_only", body["data_tiering"].(map[string]interface{})["type"])
					}
					indexState = body
					indexState["id"] = "index-id"
					indexState["default"] = false
					respond(indexState)
				case "GET /system/indices/index_sets/index-id":
					if indexState == nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					respond(indexState)
				case "PUT /system/indices/index_sets/index-id":
					assert.NotContains(t, body, "creation_date")
					assert.NotContains(t, body, "index_prefix")
					assert.NotContains(t, body, "index_analyzer")
					assert.NotContains(t, body, "default")
					assert.Equal(t, "index-id", body["id"])
					for key, value := range body {
						indexState[key] = value
					}
					respond(indexState)
				case "DELETE /system/indices/index_sets/index-id":
					indexState = nil
					w.WriteHeader(http.StatusNoContent)
				case "POST /streams":
					assert.Len(t, body, 1)
					entity, ok := body["entity"].(map[string]interface{})
					if !assert.True(t, ok) {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					assert.Len(t, entity, 5)
					assert.NotContains(t, entity, "disabled")
					streamState = entity
					streamState["id"] = "stream-id"
					streamState["disabled"] = true
					streamState["created_at"] = "2026-10-13T00:00:00Z"
					streamState["creator_user_id"] = "admin"
					streamState["is_default"] = false
					respond(map[string]interface{}{"stream_id": "stream-id"})
				case "POST /streams/stream-id/resume":
					streamState["disabled"] = false
					w.WriteHeader(http.StatusNoContent)
				case "GET /streams/stream-id":
					if streamState == nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					respond(streamState)
				case "PUT /streams/stream-id":
					assert.Len(t, body, 5)
					assert.NotContains(t, body, "entity")
					for key, value := range body {
						streamState[key] = value
					}
					respond(streamState)
				case "DELETE /streams/stream-id":
					streamState = nil
					w.WriteHeader(http.StatusNoContent)
				case "POST /streams/stream-id/rules":
					assert.Len(t, body, 5)
					assert.NotContains(t, body, "stream_id")
					assert.NotContains(t, body, "rule_id")
					ruleState = body
					ruleState["id"] = "rule-id"
					respond(map[string]interface{}{"streamrule_id": "rule-id"})
				case "GET /streams/stream-id/rules/rule-id":
					if ruleState == nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					respond(ruleState)
				case "PUT /streams/stream-id/rules/rule-id":
					assert.Len(t, body, 5)
					for key, value := range body {
						ruleState[key] = value
					}
					respond(ruleState)
				case "DELETE /streams/stream-id/rules/rule-id":
					ruleState = nil
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, path)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()
			cfg := config.Config{Endpoint: server.URL + "/api", APIVersion: "v7"}
			index := indexset.Resource()
			values := map[string]interface{}{
				"title": "test", "description": "test", "index_prefix": "test", "index_analyzer": "standard", "shards": 2, "replicas": 1,
				"writable": true, "index_optimization_max_num_segments": 1, "field_type_refresh_interval": 120000,
				"rotation_strategy_class":  "org.graylog2.indexer.rotation.strategies.SizeBasedRotationStrategy",
				"rotation_strategy":        `{"type":"org.graylog2.indexer.rotation.strategies.SizeBasedRotationStrategyConfig","max_size":10737418240}`,
				"retention_strategy_class": "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy",
				"retention_strategy":       `{"type":"org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig","max_number_of_indices":5}`,
				"use_legacy_rotation":      legacy,
			}
			if !legacy {
				values["data_tiering"] = `{"type":"hot_only","index_lifetime_min":"P14D","index_lifetime_max":"P21D"}`
			}
			id := schema.TestResourceDataRaw(t, index.Schema, values)
			require.NoError(t, index.Create(id, cfg))
			require.NoError(t, index.Read(id, cfg))
			require.Equal(t, "index-id", id.Id())
			require.NoError(t, id.Set("title", "updated"))
			require.NoError(t, index.Update(id, cfg))
			require.NoError(t, index.Read(id, cfg))
			require.Equal(t, "updated", id.Get("title"))
			st := stream.Resource()
			sd := schema.TestResourceDataRaw(t, st.Schema, map[string]interface{}{"title": "test", "description": "test", "index_set_id": id.Id(), "matching_type": "AND", "disabled": false})
			require.NoError(t, st.Create(sd, cfg))
			require.NoError(t, st.Read(sd, cfg))
			require.Equal(t, "stream-id", sd.Id())
			require.Equal(t, false, sd.Get("disabled"))
			require.NoError(t, sd.Set("title", "updated"))
			require.NoError(t, st.Update(sd, cfg))
			require.NoError(t, st.Read(sd, cfg))
			require.Equal(t, "updated", sd.Get("title"))
			rl := rule.Resource()
			rd := schema.TestResourceDataRaw(t, rl.Schema, map[string]interface{}{"stream_id": sd.Id(), "field": "source", "type": 1, "value": "test", "description": "test", "inverted": false})
			require.NoError(t, rl.Create(rd, cfg))
			require.NoError(t, rl.Read(rd, cfg))
			require.Equal(t, "stream-id/rule-id", rd.Id())
			require.NoError(t, rd.Set("value", "updated"))
			require.NoError(t, rl.Update(rd, cfg))
			require.NoError(t, rl.Read(rd, cfg))
			require.Equal(t, "updated", rd.Get("value"))
			for _, item := range []struct {
				resource *schema.Resource
				data     *schema.ResourceData
			}{{rl, rd}, {st, sd}, {index, id}} {
				require.NoError(t, item.resource.Delete(item.data, cfg))
				require.NoError(t, item.resource.Read(item.data, cfg))
				require.Empty(t, item.data.Id())
			}
			require.Contains(t, calls, "POST /streams/stream-id/resume")
		})
	}
}
