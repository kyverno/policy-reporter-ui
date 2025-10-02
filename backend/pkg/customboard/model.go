package customboard

type AccessControl struct {
	Emails []string `koanf:"emails"`
	Groups []string `koanf:"groups"`
}

type Filter struct {
	NamespaceKinds []string `koanf:"namespaceKinds"`
	ClusterKinds   []string `koanf:"clusterKinds"`
	Results        []string `koanf:"results"`
	Severities     []string `koanf:"severities"`
}

type NamespaceSelector struct {
	Selector map[string]string `koanf:"selector"`
	List     []string          `koanf:"list"`
}

type SourceSelector struct {
	List []string `koanf:"list"`
}

type PolicyReportSelector struct {
	Selector map[string]string `koanf:"selector"`
}

type FilterList struct {
	Include Filter `koanf:"include"`
	Exclude Filter `koanf:"exclude"`
}

type ClusterScope struct {
	Enabled bool `koanf:"enabled"`
}

type CustomBoard struct {
	ID            string
	Name          string               `koanf:"name"`
	AccessControl AccessControl        `koanf:"accessControl"`
	Namespaces    NamespaceSelector    `koanf:"namespaces"`
	Sources       SourceSelector       `koanf:"sources"`
	Filter        FilterList           `koanf:"filter"`
	PolicyReports PolicyReportSelector `koanf:"policyReports"`
	Display       string               `json:"display"`
	ClusterScope  ClusterScope         `koanf:"clusterScope"`
}
