package core

type Category struct {
	Name  string `json:"name"`
	Pass  int    `json:"pass"`
	Skip  int    `json:"skip"`
	Warn  int    `json:"warn"`
	Error int    `json:"error"`
	Fail  int    `json:"fail"`
}

type Policy struct {
	Source   string         `json:"source,omitempty"`
	Category string         `json:"category,omitempty"`
	Name     string         `json:"policy"`
	Severity string         `json:"severity,omitempty"`
	Results  map[string]int `json:"results"`
}

type Resource struct {
	ID         string `json:"id,omitempty"`
	UID        string `json:"uid,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type ResourceStatusCount struct {
	Source string `json:"source,omitempty"`
	Pass   int    `json:"pass"`
	Warn   int    `json:"warn"`
	Fail   int    `json:"fail"`
	Error  int    `json:"error"`
	Skip   int    `json:"skip"`
}

type SourceCategoryTree struct {
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

type FindingCounts struct {
	Total  int            `json:"total"`
	Source string         `json:"source"`
	Counts map[string]int `json:"counts"`
}

type Findings struct {
	Total     int             `json:"total"`
	PerResult map[string]int  `json:"perResult"`
	Counts    []FindingCounts `json:"counts"`
}

type NamespaceStatusCounts = map[string]map[string]int
