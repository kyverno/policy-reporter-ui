package service

import (
	"fmt"
	"sort"

	plugin "github.com/kyverno/policy-reporter-plugins/sdk/api"
	pluginAPI "github.com/kyverno/policy-reporter-plugins/sdk/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	"golang.org/x/exp/maps"
)

func MapFindingSourcesToSourceItem(findings *core.Findings) []SourceItem {
	findingSources := make(map[string]bool, 0)
	for _, f := range findings.Counts {
		findingSources[f.Source] = true
	}

	sourceItems := make([]SourceItem, 0, len(findingSources))
	for f := range findingSources {
		sourceItems = append(sourceItems, SourceItem{
			Name:  f,
			Title: utils.Title(f),
		})
	}

	sort.SliceStable(sourceItems, func(a, b int) bool {
		return sourceItems[a].Title < sourceItems[b].Title
	})

	return sourceItems
}

func MapFindingSourcesToFindingCharts(findings *core.Findings) map[string]*Chart {
	charts := make(map[string]*Chart, 0)
	totals := make(map[string]int, 0)

	for _, finding := range findings.Counts {
		for status, count := range finding.Counts {
			if chart, ok := charts[status]; ok {
				chart.Labels = append(chart.Labels, utils.Title(finding.Source))
				chart.Datasets[0].Data = append(chart.Datasets[0].Data, count)
				totals[status] += count
				continue
			}

			charts[status] = &Chart{
				Labels: []string{utils.Title(finding.Source)},
				Datasets: []*Dataset{
					{Label: status, Data: []int{count}},
				},
			}

			totals[status] = count
		}
	}

	for s, chart := range charts {
		sort.SliceStable(chart.Labels, func(i, j int) bool {
			return chart.Datasets[0].Data[i] < chart.Datasets[0].Data[j]
		})

		sort.SliceStable(chart.Datasets[0].Data, func(i, j int) bool {
			return chart.Datasets[0].Data[i] < chart.Datasets[0].Data[j]
		})

		chart.Name = fmt.Sprintf("%d", totals[s])
	}

	return charts
}

func MapFindingsToSourceStatusChart(title string, findings *core.Findings) *Chart {
	if len(findings.Counts) == 0 {
		return &Chart{
			Name:     title,
			Labels:   make([]string, 0),
			Datasets: []*Dataset{{Data: make([]int, 0, 0)}},
		}
	}

	counts := findings.Counts[0]

	values := map[string]int{
		StatusPass:  counts.Counts[StatusPass],
		StatusFail:  counts.Counts[StatusFail],
		StatusWarn:  counts.Counts[StatusWarn],
		StatusError: counts.Counts[StatusError],
		StatusSkip:  counts.Counts[StatusSkip],
	}

	labels := make([]string, 0, 5)
	dataset := &Dataset{Data: make([]int, 0, 5)}

	for s, c := range values {
		if c == 0 {
			continue
		}

		labels = append(labels, utils.Title(s))
		dataset.Data = append(dataset.Data, c)
	}

	return &Chart{
		Name:     utils.Title(counts.Source),
		Labels:   labels,
		Datasets: []*Dataset{dataset},
	}
}

func MapNamespaceStatusCountsToChart(title string, namespaces core.NamespaceStatusCounts) *Chart {
	sets := map[string]*Dataset{
		StatusPass:  {Label: utils.Title(StatusPass), Data: make([]int, 0)},
		StatusFail:  {Label: utils.Title(StatusFail), Data: make([]int, 0)},
		StatusWarn:  {Label: utils.Title(StatusWarn), Data: make([]int, 0)},
		StatusError: {Label: utils.Title(StatusError), Data: make([]int, 0)},
		StatusSkip:  {Label: utils.Title(StatusSkip), Data: make([]int, 0)},
	}

	labels := make([]string, 0, len(namespaces))
	sorting := map[string]int{}
	index := 0

	for namespace, results := range namespaces {
		labels = append(labels, namespace)
		sorting[namespace] = index
		index++

		for status, count := range results {
			sets[status].Data = append(sets[status].Data, count)
		}
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	// sorting Data to the same order as related labels
	for _, set := range sets {
		data := make([]int, 0, len(set.Data))
		for _, label := range labels {
			data = append(data, set.Data[sorting[label]])
		}

		set.Data = data
	}

	return &Chart{
		Name:   title,
		Labels: labels,
		Datasets: []*Dataset{
			sets[StatusPass],
			sets[StatusWarn],
			sets[StatusFail],
			sets[StatusError],
			sets[StatusSkip],
		},
	}
}

func MapNamespaceStatusCountsToCharts(findings map[string]core.NamespaceStatusCounts) map[string]*ChartVariants {
	charts := make(map[string]*ChartVariants, len(findings))

	for source, namespaces := range findings {
		charts[source] = &ChartVariants{
			Complete: MapNamespaceStatusCountsToChart(utils.Title(source), namespaces),
		}

		if len(namespaces) > 8 {
			prev := make(core.NamespaceStatusCounts, 8)
			ns := maps.Keys(namespaces)
			sort.Strings(ns)

			for _, v := range ns {
				if hasFailure(namespaces[v]) {
					prev[v] = namespaces[v]
				}

				if len(prev) >= 7 {
					break
				}
			}

			if len(prev) == 0 {
				for _, v := range ns[0:6] {
					prev[v] = namespaces[v]
				}
			}

			charts[source].Preview = MapNamespaceStatusCountsToChart(utils.Title(source), prev)
		}
	}

	return charts
}

func hasFailure(ns map[string]int) bool {
	var failures int
	if v, ok := ns["fail"]; ok {
		failures += v
	}
	if v, ok := ns["warn"]; ok {
		failures += v
	}
	if v, ok := ns["error"]; ok {
		failures += v
	}

	return failures > 0
}

func SumResourceCounts(results []core.ResourceStatusCount) map[string]int {
	values := map[string]int{
		StatusPass:  0,
		StatusFail:  0,
		StatusWarn:  0,
		StatusError: 0,
		StatusSkip:  0,
	}

	for _, result := range results {
		values[StatusPass] += result.Pass
		values[StatusWarn] += result.Warn
		values[StatusFail] += result.Fail
		values[StatusError] += result.Error
		values[StatusSkip] += result.Skip
	}

	return values
}

func MapResourceSourceChart(results []core.ResourceStatusCount) *Chart {
	sets := map[string]*Dataset{
		StatusPass:  {Label: utils.Title(StatusPass), Data: make([]int, 0)},
		StatusFail:  {Label: utils.Title(StatusFail), Data: make([]int, 0)},
		StatusWarn:  {Label: utils.Title(StatusWarn), Data: make([]int, 0)},
		StatusError: {Label: utils.Title(StatusError), Data: make([]int, 0)},
		StatusSkip:  {Label: utils.Title(StatusSkip), Data: make([]int, 0)},
	}

	labels := make([]string, 0, len(results))
	sorting := make(map[string]int, len(results))
	for index, result := range results {
		label := utils.Title(result.Source)
		sorting[label] = index
		labels = append(labels, label)

		sets[StatusPass].Data = append(sets[StatusPass].Data, result.Pass)
		sets[StatusWarn].Data = append(sets[StatusWarn].Data, result.Warn)
		sets[StatusFail].Data = append(sets[StatusFail].Data, result.Fail)
		sets[StatusError].Data = append(sets[StatusError].Data, result.Error)
		sets[StatusSkip].Data = append(sets[StatusSkip].Data, result.Skip)
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	// sorting Data to the same order as related labels
	for _, set := range sets {
		data := make([]int, 0, len(set.Data))
		for _, label := range labels {
			data = append(data, set.Data[sorting[label]])
		}

		set.Data = data
	}

	return &Chart{
		Labels: labels,
		Datasets: []*Dataset{
			sets[StatusPass],
			sets[StatusWarn],
			sets[StatusFail],
			sets[StatusError],
			sets[StatusSkip],
		},
	}
}

func MapCategoriesToChart(title string, categories []core.Category) *Chart {
	sets := map[string]*Dataset{
		StatusPass:  {Label: utils.Title(StatusPass), Data: make([]int, 0)},
		StatusFail:  {Label: utils.Title(StatusFail), Data: make([]int, 0)},
		StatusWarn:  {Label: utils.Title(StatusWarn), Data: make([]int, 0)},
		StatusError: {Label: utils.Title(StatusError), Data: make([]int, 0)},
		StatusSkip:  {Label: utils.Title(StatusSkip), Data: make([]int, 0)},
	}

	labels := make([]string, 0, len(categories))
	sorting := make(map[string]int, len(categories))

	for index, category := range categories {
		sorting[category.Name] = index
		labels = append(labels, category.Name)

		sets[StatusPass].Data = append(sets[StatusPass].Data, category.Pass)
		sets[StatusWarn].Data = append(sets[StatusWarn].Data, category.Warn)
		sets[StatusFail].Data = append(sets[StatusFail].Data, category.Fail)
		sets[StatusError].Data = append(sets[StatusError].Data, category.Error)
		sets[StatusSkip].Data = append(sets[StatusSkip].Data, category.Skip)
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	// sorting Data to the same order as related labels
	for _, set := range sets {
		data := make([]int, 0, len(set.Data))
		for _, label := range labels {
			data = append(data, set.Data[sorting[label]])
		}

		set.Data = data
	}

	return &Chart{
		Name:   title,
		Labels: labels,
		Datasets: []*Dataset{
			sets[StatusPass],
			sets[StatusWarn],
			sets[StatusFail],
			sets[StatusError],
			sets[StatusSkip],
		},
	}
}

func MapPolicyDetails(details *PolicyDetails, policy *plugin.Policy) *PolicyDetails {
	if policy == nil {
		return details
	}

	details.Title = policy.Title
	details.Description = policy.Description
	details.Severity = policy.Severity
	details.References = utils.Map(policy.References, func(ref pluginAPI.Reference) string {
		return ref.URL
	})

	details.ShowDetails = true

	if policy.Engine != nil {
		details.Engine = &Engine{
			Name:     policy.Engine.Name,
			Version:  policy.Engine.Version,
			Subjects: policy.Engine.Subjects,
		}
	}

	if policy.SourceCode != nil {
		details.SourceCode = &SourceCode{
			ContentType: policy.SourceCode.ContentType,
			Content:     policy.SourceCode.Content,
		}
	}

	details.Additional = utils.Map(policy.Additional, func(d pluginAPI.Details) Details {
		return Details{
			Title: d.Title,
			Items: utils.Map(d.Items, func(i pluginAPI.DetailsItem) Item {
				return Item{Title: i.Title, Value: i.Value}
			}),
		}
	})

	details.Details = utils.Map(policy.Details, func(i pluginAPI.DetailsItem) Item {
		return Item{Title: i.Title, Value: i.Value}
	})

	return details
}
