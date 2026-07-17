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

	NamespaceKinds ValueFilter `koanf:"namespaceKinds"`
	ClusterKinds   ValueFilter `koanf:"clusterKinds"`
	Results        ValueFilter `koanf:"results"`
	Severities     ValueFilter `koanf:"severities"`
}

type ClusterScope struct {
	Enabled bool `koanf:"enabled"`
}

type RenderOptions struct {
	DashboardMode string `koanf:"dashboardMode"`
	ResultView    string `koanf:"resultView"`
}

type CustomBoard struct {
	ID              string
	Name            string               `koanf:"name"`
	AccessControl   AccessControl        `koanf:"accessControl"`
	Namespaces      NamespaceSelector    `koanf:"namespaces"`
	Sources         SourceSelector       `koanf:"sources"`
	Filter          FilterList           `koanf:"filter"`
	PolicyReports   PolicyReportSelector `koanf:"policyReports"`
	Display         string               `koanf:"display"`
	AllowedDisplays []string             `koanf:"allowedDisplays"`
	RenderOptions   RenderOptions        `koanf:"renderOptions"`
	ClusterScope    ClusterScope         `koanf:"clusterScope"`
}

func (c *CustomBoard) AllowedResultViews() []string {
	if len(c.AllowedDisplays) > 0 {
		return c.AllowedDisplays
	}

	if c.Display != "" {
		return []string{c.Display}
	}
	if c.RenderOptions.ResultView != "" {
		return []string{c.RenderOptions.ResultView}
	}

	return []string{"resources"}
}

func (c *CustomBoard) ResultView() string {
	views := c.AllowedResultViews()
	for _, view := range views {
		if view == c.Display {
			return c.Display
		}
	}

	return views[0]
}
