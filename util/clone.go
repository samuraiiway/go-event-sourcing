package util

func CloneMap(original map[string]interface{}) map[string]interface{} {
	clone := map[string]interface{}{}

	for key, value := range original {
		clone[key] = value
	}

	return clone
}
