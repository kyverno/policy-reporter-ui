package utils

func Contains[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

func Some[T comparable](list []T, items []T) bool {
	for _, i := range list {
		for _, j := range items {
			if i == j {
				return true
			}
		}
	}

	return false
}

func Intersect[T comparable](list []T, subset []T) []T {
	compare := ToCompareMap(list)

	result := make([]T, 0, len(subset))
	for _, i := range subset {
		if _, ok := compare[i]; ok {
			result = append(result, i)
		}
	}

	return result
}

func ToCompareMap[T comparable](list []T) map[T]struct{} {
	result := make(map[T]struct{}, len(list))
	for _, i := range list {
		result[i] = struct{}{}
	}

	return result
}
