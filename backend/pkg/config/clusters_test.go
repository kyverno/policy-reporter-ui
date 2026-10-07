package config

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestLoadCustomBoardClusters(t *testing.T) {
	t.Parallel()
	filename := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(filename, []byte(`customBoards:
  - name: Restricted
    clusters: [production, staging]
  - name: Unrestricted
  - name: Empty
    clusters: []
`), 0o600); err != nil {
		t.Fatal(err)
	}
	var config Config
	if err := Load(&config, filename); err != nil {
		t.Fatal(err)
	}
	if len(config.CustomBoards) != 3 {
		t.Fatalf("boards = %d", len(config.CustomBoards))
	}
	if !slices.Equal(config.CustomBoards[0].Clusters, []string{"production", "staging"}) {
		t.Fatalf("clusters = %v", config.CustomBoards[0].Clusters)
	}
	for _, board := range config.CustomBoards[1:] {
		if !board.MatchesCluster("anything") {
			t.Fatalf("%s must be unrestricted", board.Name)
		}
	}
}
