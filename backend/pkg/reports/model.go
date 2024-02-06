package reports

import (
	"slices"

	"github.com/kyverno/policy-reporter-plugins/sdk/api"
)

const (
	StatusPass  = "pass"
	StatusFail  = "fail"
	StatusWarn  = "warn"
	StatusError = "error"
	StatusSkip  = "skip"
)

type Summary struct {
	Error   int
	Pass    int
	Fail    int
	Warning int
}

type Resource struct {
	Kind       string
	APIVersion string
	Name       string
	Status     string
}

type Rule struct {
	Summary   *Summary
	Resources []*Resource
}

type Group struct {
	Name    string
	Policy  api.PolicyListItem
	Summary *Summary
	Rules   map[string]*Rule
}

type Policy struct {
	Title       string
	Category    string
	Description string
	Severity    string
}

type Validation struct {
	Name   string
	Policy api.PolicyListItem
	Groups map[string]*Group
}

type Filter struct {
	Namespaces   []string
	Policies     []string
	Categories   []string
	Kinds        []string
	ClusterScope bool
}

func (f Filter) IncludesPolicy(policy string) bool {
	if len(f.Policies) == 0 {
		return true
	}

	return slices.Contains(f.Policies, policy)
}

func (f Filter) IncludesNamespace(namespace string) bool {
	if len(f.Namespaces) == 0 || namespace == "" {
		return true
	}

	return slices.Contains(f.Namespaces, namespace)
}
