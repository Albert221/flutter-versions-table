package utils

func MapSlice[T, S any](source []T, mapper func(T) S) []S {
	result := make([]S, len(source))
	for i, v := range source {
		result[i] = mapper(v)
	}
	return result
}

func ReverseSlice[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
