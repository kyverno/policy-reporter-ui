package utils

import (
	"net/url"
	"strings"

	"go.uber.org/zap"
)

func BasePath(callback string) string {
	if strings.HasSuffix(callback, "/callback") {
		callback = strings.TrimSuffix(callback, "/callback")
	}

	r, err := url.Parse(callback)
	if err != nil {
		zap.L().Error("failed to parse URL", zap.String("url", callback))

		return "/"
	}

	if r.Path == "" {
		return "/"
	}

	if !strings.HasPrefix(r.Path, "/") {
		r.Path = "/" + r.Path
	}

	if strings.HasSuffix(r.Path, "/") {
		return r.Path
	}

	return r.Path + "/"
}
