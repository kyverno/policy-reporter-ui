package customboard

import (
	"github.com/gosimple/slug"
	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/customboard/v1alpha1"
)

func MapCustomBoardToModel(cb *v1alpha1.CustomBoard) *CustomBoard {
	return &CustomBoard{
		ID:            slug.Make(cb.Spec.Title),
		Name:          cb.Spec.Title,
		AccessControl: MapAccessControl(cb.Spec.AccessControl),
		Filter: FilterList{
			Include: MapIncludeFilter(cb.Spec.Filter),
			Exclude: MapExcludeFilter(cb.Spec.Filter),
		},
		Display:       string(cb.Spec.Display),
		Namespaces:    NamespaceSelector{Selector: cb.Spec.NamespaceSelector.LabelSelector, List: cb.Spec.NamespaceSelector.List},
		Sources:       MapSources(cb.Spec.SourceSelector),
		PolicyReports: MapPolicyReports(cb.Spec.PolicyReportSelector),
		ClusterScope:  ClusterScope{Enabled: cb.Spec.ClusterScope == nil || cb.Spec.ClusterScope.Enabled},
	}
}

func MapNamespaceCustomBoardToModel(cb *v1alpha1.NamespaceCustomBoard) *CustomBoard {
	return &CustomBoard{
		ID:            slug.Make(cb.Spec.Title),
		Name:          cb.Spec.Title,
		AccessControl: MapAccessControl(cb.Spec.AccessControl),
		Filter: FilterList{
			Include: MapIncludeFilter(cb.Spec.Filter),
			Exclude: MapExcludeFilter(cb.Spec.Filter),
		},
		Display:       string(cb.Spec.Display),
		Namespaces:    NamespaceSelector{List: []string{cb.Namespace}},
		Sources:       MapSources(cb.Spec.SourceSelector),
		PolicyReports: MapPolicyReports(cb.Spec.PolicyReportSelector),
		ClusterScope:  ClusterScope{Enabled: false},
	}
}

func MapAccessControl(ac *v1alpha1.AccessControl) AccessControl {
	if ac == nil {
		return AccessControl{}
	}

	return AccessControl{
		Emails: ac.Emails,
		Groups: ac.Groups,
	}
}

func MapIncludeFilter(f *v1alpha1.Filter) Filter {
	if f == nil {
		return Filter{}
	}

	return Filter{
		NamespaceKinds: f.NamespaceKinds.Include,
		ClusterKinds:   f.ClusterKinds.Include,
		Results:        f.Results.Include,
		Severities:     f.Severities.Include,
	}
}

func MapExcludeFilter(f *v1alpha1.Filter) Filter {
	if f == nil {
		return Filter{}
	}

	return Filter{
		NamespaceKinds: f.NamespaceKinds.Exclude,
		ClusterKinds:   f.ClusterKinds.Exclude,
		Results:        f.Results.Exclude,
		Severities:     f.Severities.Exclude,
	}
}

func MapSources(s *v1alpha1.SourceSelector) SourceSelector {
	if s == nil {
		return SourceSelector{}
	}

	return SourceSelector{
		List: s.List,
	}
}

func MapPolicyReports(pr *v1alpha1.PolicyReportSelector) PolicyReportSelector {
	if pr == nil {
		return PolicyReportSelector{}
	}

	return PolicyReportSelector{
		Selector: pr.LabelSelector,
	}
}
