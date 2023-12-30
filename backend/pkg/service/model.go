package service

import core "github.com/kyverno/policy-reporter-ui/pkg/core/client"

const (
	StatusPass  = "pass"
	StatusFail  = "fail"
	StatusWarn  = "warn"
	StatusError = "error"
	StatusSkip  = "skip"
)

type SourceItem struct {
	Title string `json:"title"`
	Name  string `json:"name"`
}

type Dataset struct {
	Label           string `json:"label"`
	BackgroundColor string `json:"backgroundColor"`
	Data            []int  `json:"data"`
}

type Chart struct {
	Labels   []string   `json:"labels"`
	Datasets []*Dataset `json:"datasets"`
	Name     string     `json:"name"`
}

type Charts struct {
	ClusterScope   map[string]map[string]int `json:"clusterScope"`
	NamespaceScope map[string]*Chart         `json:"namespaceScope"`
	Findings       any                       `json:"findings"`
}

type Total struct {
	Count     int            `json:"count"`
	PerResult map[string]int `json:"perResult"`
}

type Dashboard struct {
	FilterSources  []string     `json:"filterSources,omitempty"`
	ClusterScope   bool         `json:"clusterScope"`
	Sources        []string     `json:"sources"`
	Namespaces     []string     `json:"namespaces"`
	SingleSource   bool         `json:"singleSource"`
	MultipleSource bool         `json:"multiSource"`
	Charts         Charts       `json:"charts"`
	SourcesNavi    []SourceItem `json:"sourcesNavi"`
	Total          Total        `json:"total"`
}

type ResourceDetails struct {
	Resource *core.Resource `json:"resource"`
	Results  map[string]int `json:"results"`
	Chart    *Chart         `json:"chart"`
	Sources  []Source       `json:"sources"`
}

type Source struct {
	Title      string   `json:"title"`
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
	Chart      *Chart   `json:"chart"`
}

type PolicyCharts struct {
	Findings       *Chart         `json:"findings"`
	NamespaceScope *Chart         `json:"namespaceScope"`
	ClusterScope   map[string]int `json:"clusterScope"`
}

type PolicyDetails struct {
	Title      string       `json:"title"`
	Name       string       `json:"name"`
	Namespaces []string     `json:"namespaces"`
	Chart      PolicyCharts `json:"charts"`
}
