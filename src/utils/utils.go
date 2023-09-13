package utils

func GetMapKeys[K comparable, V any](value map[K]V) []K {
	result := []K{}
	for key := range value {
		result = append(result, key)
	}
	return result
}
