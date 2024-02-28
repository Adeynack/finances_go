package utils

func MapGetValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, len(m))
	var index uint
	for _, value := range m {
		values[index] = value
		index += 1
	}
	return values
}
