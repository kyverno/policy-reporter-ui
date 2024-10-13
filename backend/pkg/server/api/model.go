package api

import "github.com/kyverno/policy-reporter-ui/pkg/auth"

type Policy struct {
	Source      string         `json:"source,omitempty"`
	Category    string         `json:"category,omitempty"`
	Namespace   string         `json:"namespace,omitempty"`
	Name        string         `json:"name"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Severity    string         `json:"severity,omitempty"`
	Results     map[string]int `json:"results"`
}

type DefaultFilter struct {
	Resources        []string `json:"resources"`
	ClusterResources []string `json:"clusterResources"`
}

type Excludes struct {
	NamespaceKinds []string `json:"namespaceKinds"`
	ClusterKinds   []string `json:"clusterKinds"`
	Results        []string `json:"results"`
	Severities     []string `json:"severities"`
}

type Source struct {
	Name       string   `json:"name"`
	ViewType   string   `mapstructure:"type"`
	Exceptions bool     `mapstructure:"exceptions"`
	Excludes   Excludes `json:"excludes"`
}

type Cluster struct {
	auth.Permissions `json:"-"`
	Name             string   `json:"name"`
	Slug             string   `json:"slug"`
	Plugins          []string `json:"plugins"`
}

type PolicyReports struct {
	Selector map[string]string
}

type Namespaces struct {
	Selector map[string]string
	List     []string
}

type Sources struct {
	List []string
}

type CustomBoard struct {
	auth.Permissions `json:"-"`
	Name             string        `json:"name"`
	ID               string        `json:"id"`
	ClusterScope     bool          `json:"-"`
	Namespaces       Namespaces    `json:"-"`
	Sources          Sources       `json:"-"`
	PolicyReports    PolicyReports `json:"-"`
}

type Boards struct {
	auth.Permissions
}

type Config struct {
	User     any       `json:"user"`
	Clusters []Cluster `json:"clusters"`
	Sources  []Source  `json:"sources"`
	Default  string    `json:"default"`
	Boards   Boards    `json:"-"`
	Banner   string    `json:"banner"`
	OAuth    bool      `json:"oauth"`
}

type NavigationItem struct {
	Title    string           `json:"title"`
	Subtitle string           `json:"subtitle"`
	Path     string           `json:"path"`
	Children []NavigationItem `json:"children,omitempty"`
}
