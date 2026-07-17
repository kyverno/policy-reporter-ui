package customboard

type AccessControl struct {
	Emails []string `koanf:"emails"`
	Groups []string `koanf:"groups"`
}

type ValueFilter struct {
	Include []string `koanf:"include"`
	Exclude []string `koanf:"exclude"`
}

type Filter struct {
	NamespaceKinds   []string `koanf:"namespaceKinds"`
	ClusterKinds     []string `koanf:"clusterKinds"`
	Results          []string `koanf:"results"`
	Severities       []string `koanf:"severities"`
	Resources        []string `koanf:"resources"`
	ClusterResources []string `koanf:"clusterResources"`
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

	NamespaceKinds   ValueFilter `koanf:"namespaceKinds"`
	ClusterKinds     ValueFilter `koanf:"clusterKinds"`
	Results          ValueFilter `koanf:"results"`
	Severities       ValueFilter `koanf:"severities"`
	Resources        ValueFilter `koanf:"resources"`
	ClusterResources ValueFilter `koanf:"clusterResources"`
}

type ClusterScope struct {
	Enabled bool `koanf:"enabled"`
}

type RenderOptions struct {
	DashboardMode string `koanf:"dashboardMode"`
	ResultView    string `koanf:"resultView"`
}

type CustomBoard struct {
	ID            string
	Name          string               `koanf:"name"`
	AccessControl AccessControl        `koanf:"accessControl"`
	Namespaces    NamespaceSelector    `koanf:"namespaces"`
	Sources       SourceSelector       `koanf:"sources"`
	Filter        FilterList           `koanf:"filter"`
	PolicyReports PolicyReportSelector `koanf:"policyReports"`
	Display       string               `koanf:"display"`
	RenderOptions RenderOptions        `koanf:"renderOptions"`
	ClusterScope  ClusterScope         `koanf:"clusterScope"`
}

func (c *CustomBoard) ResultView() string {
	if c.RenderOptions.ResultView != "" {
		return c.RenderOptions.ResultView
	}
	return c.Display
}
