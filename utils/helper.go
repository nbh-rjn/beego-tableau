package utils

func DatasourceExists(datasources []map[string]interface{}, nameToFind string) (string, bool) {
	for _, datasource := range datasources {
		if name, ok := datasource["name"].(string); ok && name == nameToFind {
			return datasource["id"].(string), true
		}
	}
	return "", false
}
