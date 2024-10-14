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
