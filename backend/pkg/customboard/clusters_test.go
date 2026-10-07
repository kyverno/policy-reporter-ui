package customboard

import (
	"slices"
	"testing"

	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/ui/v1alpha1"
)

func TestMatchesCluster(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name     string
		clusters []string
		want     bool
	}{
		{"nil", nil, true},
		{"empty", []string{}, true},
		{"match", []string{"production"}, true},
		{"nonmatch", []string{"staging"}, false},
		{"multiple", []string{"staging", "production"}, true},
		{"duplicates", []string{"production", "production"}, true},
		{"unknown", []string{"unknown"}, false},
		{"case-sensitive", []string{"Production"}, false},
		{"display-name-not-slug", []string{"Production Cluster"}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			board := CustomBoard{Clusters: tt.clusters}
			if got := board.MatchesCluster("production"); got != tt.want {
				t.Fatalf("MatchesCluster = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapCustomBoardClusters(t *testing.T) {
	t.Parallel()
	for _, clusters := range [][]string{nil, {}, {"production", "staging"}} {
		board := &v1alpha1.CustomBoard{Spec: v1alpha1.CustomBoardSpec{Clusters: clusters}}
		mapped := MapCustomBoardToModel(board)
		if !slices.Equal(mapped.Clusters, clusters) {
			t.Fatalf("clusters = %v, want %v", mapped.Clusters, clusters)
		}
		if len(clusters) == 0 && !mapped.MatchesCluster("anything") {
			t.Fatal("omitted clusters must be unrestricted")
		}
	}
}

func TestCustomBoardDeepCopyClusters(t *testing.T) {
	t.Parallel()
	board := &v1alpha1.CustomBoard{Spec: v1alpha1.CustomBoardSpec{Clusters: []string{"production"}}}
	copied := board.DeepCopy()
	copied.Spec.Clusters[0] = "staging"
	if board.Spec.Clusters[0] != "production" {
		t.Fatal("deep copy aliases clusters")
	}
}
