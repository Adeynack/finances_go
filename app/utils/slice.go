package utils

func ToStringSlice[T ~string](list []T) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = string(v)
	}
	return result
}
