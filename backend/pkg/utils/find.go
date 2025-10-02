package utils

func Find[T any](s []T, keep func(T) bool, fallback T) T {
	for _, n := range s {
		if keep(n) {
			return n
		}
	}
	return fallback
}

func FindInMap[R comparable, T any](s map[R]T, keep func(T) bool, fallback T) T {
	for _, n := range s {
		if keep(n) {
			return n
		}
	}
	return fallback
}
