package utils

func CloneMap[R comparable, K any, T ~map[R]K](src T) T {
	dst := make(T, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
