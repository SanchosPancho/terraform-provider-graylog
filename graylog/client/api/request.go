package api

// UsesUserIDs preserves the v4 user API behavior when using Graylog 7.
func UsesUserIDs(version string) bool {
	return version == "v4" || version == "v7"
}

// CreateEntity wraps creation requests for the entity-sharing API introduced
// in Graylog 7. Updates and index-set requests retain their flat payloads.
func CreateEntity(version string, data map[string]interface{}) map[string]interface{} {
	if version == "v7" {
		return map[string]interface{}{"entity": data}
	}
	return data
}

// Fields copies only fields supported by a request model. Response-only fields
// must not be sent back to Graylog 7's strict JSON deserializer.
func Fields(data map[string]interface{}, names ...string) map[string]interface{} {
	result := make(map[string]interface{}, len(names))
	for _, name := range names {
		if value, ok := data[name]; ok {
			result[name] = value
		}
	}
	return result
}
