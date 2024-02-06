package utils

func Filter[T any](source []T, cb func(T) bool) []T {
	list := make([]T, 0, len(source))
	for _, i := range source {
		if cb(i) {
			list = append(list, i)
		}
	}

	return list
}
