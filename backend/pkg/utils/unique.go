package utils

func Unique[T comparable](list []T) []T {
	uniqueMap := make(map[T]struct{})
	for _, item := range list {
		uniqueMap[item] = struct{}{}
	}

	uniqueList := make([]T, 0, len(uniqueMap))
	for item := range uniqueMap {
		uniqueList = append(uniqueList, item)
	}

	return uniqueList
}
