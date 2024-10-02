package api

import "slices"

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
	Name    string   `json:"name"`
	Slug    string   `json:"slug"`
	Plugins []string `json:"plugins"`
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

type AccessControl struct {
	Emails []string
}

type Permissions struct {
	AccessControl AccessControl `json:"-"`
}

func (p Permissions) AllowedEmail(email string) bool {
	if len(p.AccessControl.Emails) == 0 {
		return true
	}

	return slices.Contains(p.AccessControl.Emails, email)
}

type CustomBoard struct {
	Permissions
	Name          string        `json:"name"`
	ID            string        `json:"id"`
	ClusterScope  bool          `json:"-"`
	Namespaces    Namespaces    `json:"-"`
	Sources       Sources       `json:"-"`
	PolicyReports PolicyReports `json:"-"`
}

type Boards struct {
	Permissions
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
