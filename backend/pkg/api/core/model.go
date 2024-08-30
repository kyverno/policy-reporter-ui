package core

type Status struct {
	Pass  int `json:"pass"`
	Skip  int `json:"skip"`
	Warn  int `json:"warn"`
	Error int `json:"error"`
	Fail  int `json:"fail"`
}

type Severities struct {
	Unknown  int `json:"unknown"`
	Low      int `json:"low"`
	Info     int `json:"info"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

type Category struct {
	Name       string     `json:"name"`
	Status     Status     `json:"status"`
	Severities Severities `json:"severities"`
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

type ResourceSeverityCount struct {
	Source   string `json:"source,omitempty"`
	Info     int    `json:"info"`
	Low      int    `json:"low"`
	Medium   int    `json:"medium"`
	High     int    `json:"high"`
	Critical int    `json:"critical"`
	Unknown  int    `json:"unknown"`
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

type PolicyResult struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace,omitempty"`
	Kind       string            `json:"kind"`
	APIVersion string            `json:"apiVersion"`
	Name       string            `json:"name"`
	Message    string            `json:"message"`
	Category   string            `json:"category,omitempty"`
	Policy     string            `json:"policy"`
	Rule       string            `json:"rule"`
	Status     string            `json:"status"`
	Severity   string            `json:"severity,omitempty"`
	Timestamp  int64             `json:"timestamp,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type Paginated[T any] struct {
	Items []T `json:"items"`
	Count int `json:"count"`
}

type Target struct {
	Name string `json:"name"`
	Host string `json:"host"`
}
