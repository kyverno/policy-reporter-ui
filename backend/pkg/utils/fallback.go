package utils

func Fallback[T comparable](s, f T) T {
	var zero T
	if s != zero {
		return s
	}

	return f
}
