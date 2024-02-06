package utils

func Contains[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}
