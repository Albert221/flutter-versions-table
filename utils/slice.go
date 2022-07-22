package utils

func MapSlice[T, S any](source []T, mapper func(T) S) []S {
	result := make([]S, len(source))
	for i, v := range source {
		result[i] = mapper(v)
	}
	return result
}
