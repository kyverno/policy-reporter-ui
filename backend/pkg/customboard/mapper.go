package customboard

import (
	"fmt"

	"github.com/gosimple/slug"

	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/customboard/v1alpha1"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

func MapCustomBoardToModel(cb *v1alpha1.CustomBoard) *CustomBoard {
	return &CustomBoard{
		ID:            slug.Make(cb.Spec.Title),
		Name:          cb.Spec.Title,
		AccessControl: MapAccessControl(cb.Spec.AccessControl),
		Filter: FilterList{
			Include: MapIncludeFilter(cb.Spec.Filter),
		},
		Display: string(cb.Spec.Display),
		RenderOptions: RenderOptions{
			ResultView:    string(utils.Fallback(cb.Spec.Display, cb.Spec.RenderOptions.ResultView)),
			DashboardMode: cb.Spec.RenderOptions.DashboardMode,
		},
		Namespaces:    NamespaceSelector{Selector: cb.Spec.NamespaceSelector.LabelSelector, List: cb.Spec.NamespaceSelector.List},
		Sources:       MapSources(cb.Spec.SourceSelector),
		PolicyReports: MapPolicyReports(cb.Spec.PolicyReportSelector),
		ClusterScope:  ClusterScope{Enabled: cb.Spec.ClusterScope == nil || cb.Spec.ClusterScope.Enabled},
	}
}

func MapNamespaceCustomBoardToModel(cb *v1alpha1.NamespaceCustomBoard) *CustomBoard {
	return &CustomBoard{
		ID:            fmt.Sprintf("%s-%s", cb.Namespace, cb.Name),
		Name:          cb.Spec.Title,
		AccessControl: MapAccessControl(cb.Spec.AccessControl),
		Filter: FilterList{
			Include: MapIncludeFilter(cb.Spec.Filter),
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
		NamespaceKinds:   f.NamespaceKinds.Include,
		ClusterKinds:     f.ClusterKinds.Include,
		Results:          f.Results.Include,
		Severities:       f.Severities.Include,
		Resources:        f.Resources.Include,
		ClusterResources: f.ClusterResources.Include,
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

func MapFilterFields(f FilterList) FilterList {
	f.Include.ClusterKinds = append(f.Include.ClusterKinds, f.ClusterKinds.Include...)
	f.Include.NamespaceKinds = append(f.Include.NamespaceKinds, f.NamespaceKinds.Include...)
	f.Include.Results = append(f.Include.Results, f.Results.Include...)
	f.Include.Severities = append(f.Include.Severities, f.Severities.Include...)
	f.Include.Resources = append(f.Include.Resources, f.Resources.Include...)
	f.Include.ClusterResources = append(f.Include.ClusterResources, f.ClusterResources.Include...)
	return f
}
