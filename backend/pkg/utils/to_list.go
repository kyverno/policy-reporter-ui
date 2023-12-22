package utils

func ToList[T any, R comparable](mapping map[R]T) []T {
	list := make([]T, 0, len(mapping))
	for _, i := range mapping {
		list = append(list, i)
	}

	return list
}
