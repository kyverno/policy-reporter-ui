package utils

func Fallback(s, f string) string {
	if s != "" {
		return s
	}

	return f
}
