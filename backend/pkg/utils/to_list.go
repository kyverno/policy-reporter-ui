package utils

func ToList[T any, R comparable](mapping map[R]T) []T {
	list := make([]T, 0, len(mapping))
	for _, i := range mapping {
		list = append(list, i)
	}

	return list
}

func Keys[T any, R comparable](mapping map[R]T) []R {
	list := make([]R, 0, len(mapping))
	for i := range mapping {
		list = append(list, i)
	}

	return list
}
