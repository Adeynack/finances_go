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

func MapGetOr[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	value, ok := m[key]
	if ok {
		return value
	}
	return defaultValue
}
