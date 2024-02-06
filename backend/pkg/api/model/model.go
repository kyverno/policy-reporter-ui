package model

import (
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
)

type Endpoints struct {
	Core    *core.Client
	Plugins map[string]*plugin.Client
}
