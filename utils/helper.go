package utils

func DatasourceExists(datasources []map[string]interface{}, nameToFind string) string {
	for _, datasource := range datasources {
		if name, ok := datasource["name"].(string); ok && name == nameToFind {
			return datasource["id"].(string)
		}
	}
	return ""
}
