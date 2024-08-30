package model

import (
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

const (
	Pass  string = "pass"
	Fail  string = "fail"
	Warn  string = "warn"
	Error string = "error"
	Skip  string = "skip"

	Unknown  string = "unknown"
	Info     string = "info"
	Low      string = "low"
	Medium   string = "medium"
	High     string = "high"
	Critical string = "critical"
)

const (
	Severity string = "severity"
	Status   string = "status"
)

type Endpoints struct {
	Core    *core.Client
	Plugins map[string]*plugin.Client
}

type SourceConfig struct {
	Results    []string
	Severities []string
	Exceptions bool
	ViewType   string
}

func (s SourceConfig) EnabledResults() []string {
	list := []string{Skip, Pass, Warn, Fail, Error}
	if len(s.Results) == 0 {
		return list
	}

	return utils.Filter(list, func(result string) bool {
		return !utils.Contains(s.Results, result)
	})
}

func (s SourceConfig) EnabledSeverities() []string {
	list := []string{Unknown, Info, Low, Medium, High, Critical}
	if len(s.Severities) == 0 {
		return list
	}

	return utils.Filter(list, func(result string) bool {
		return !utils.Contains(s.Severities, result)
	})
}
