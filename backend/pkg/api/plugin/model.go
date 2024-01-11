package plugin

import "fmt"

type Policy struct {
	Category    string `json:"category"`
	Namespace   string `json:"namespace,omitempty"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity,omitempty"`
}

func (p Policy) ID() string {
	if p.Namespace == "" {
		return p.Name
	}

	return fmt.Sprintf("%s/%s", p.Namespace, p.Name)
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
	Category    string     `json:"category"`
	Namespace   string     `json:"namespace,omitempty"`
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    string     `json:"severity,omitempty"`
	Engine      Engine     `json:"engine,omitempty"`
	SourceCode  SourceCode `json:"code,omitempty"`
	Additional  []Item     `json:"additional"`
	Details     []Details  `json:"details,omitempty"`
	References  []string   `json:"references,omitempty"`
}
