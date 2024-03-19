package service

import "github.com/kyverno/policy-reporter-ui/pkg/api/core"

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

type ChartVariants struct {
	Preview  *Chart `json:"preview"`
	Complete *Chart `json:"complete"`
}

type Charts struct {
	ClusterScope   map[string]map[string]int `json:"clusterScope"`
	NamespaceScope map[string]*ChartVariants `json:"namespaceScope"`
	Findings       any                       `json:"findings"`
}

type Total struct {
	Count     int            `json:"count"`
	PerResult map[string]int `json:"perResult"`
}

type Dashboard struct {
	Title          string       `json:"title"`
	FilterSources  []string     `json:"filterSources,omitempty"`
	ClusterScope   bool         `json:"clusterScope"`
	Sources        []string     `json:"sources"`
	Namespaces     []string     `json:"namespaces"`
	SingleSource   bool         `json:"singleSource"`
	MultipleSource bool         `json:"multiSource"`
	Charts         Charts       `json:"charts"`
	SourcesNavi    []SourceItem `json:"sourcesNavi"`
	Total          Total        `json:"total"`
	ShowResults    []string     `json:"showResults"`
	Status         []string     `json:"status"`
}

type ResourceDetails struct {
	Resource *core.Resource `json:"resource"`
	Results  map[string]int `json:"results"`
	Chart    *Chart         `json:"chart"`
	Sources  []Source       `json:"sources"`
	Status   []string       `json:"status"`
}

type Source struct {
	Title      string   `json:"title"`
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
	Status     []string `json:"status"`
	Chart      *Chart   `json:"chart"`
}

type PolicyCharts struct {
	Findings       *Chart         `json:"findings"`
	NamespaceScope *ChartVariants `json:"namespaceScope"`
	ClusterScope   map[string]int `json:"clusterScope"`
}

type Engine struct {
	Name              string   `json:"name"`
	KubernetesVersion string   `json:"kubernetesVersion,omitempty"`
	Version           string   `json:"version,omitempty"`
	Subjects          []string `json:"subjects,omitempty"`
}

type SourceCode struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type Item struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type Details struct {
	Title string `json:"title"`
	Items []Item `json:"items"`
}

type PolicyDetails struct {
	Title       string       `json:"title"`
	Name        string       `json:"name"`
	Namespaces  []string     `json:"namespaces"`
	Chart       PolicyCharts `json:"charts"`
	Description string       `json:"description"`
	Severity    string       `json:"severity,omitempty"`
	Engine      *Engine      `json:"engine,omitempty"`
	SourceCode  *SourceCode  `json:"sourceCode,omitempty"`
	Details     []Item       `json:"details,omitempty"`
	Additional  []Details    `json:"additional,omitempty"`
	References  []string     `json:"references,omitempty"`
	ShowDetails bool         `json:"showDetails"`
}
