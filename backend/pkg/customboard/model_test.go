package customboard

import "testing"

func TestCustomBoardResultViews(t *testing.T) {
	tests := []struct {
		name        string
		board       CustomBoard
		allowed     []string
		defaultView string
	}{
		{name: "legacy display", board: CustomBoard{Display: "results"}, allowed: []string{"results"}, defaultView: "results"},
		{name: "allowed displays preserve order", board: CustomBoard{Display: "results", AllowedDisplays: []string{"resources", "results"}}, allowed: []string{"resources", "results"}, defaultView: "results"},
		{name: "invalid configured default falls back", board: CustomBoard{Display: "unknown", AllowedDisplays: []string{"resources", "results"}}, allowed: []string{"resources", "results"}, defaultView: "resources"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.board.AllowedResultViews()
			if len(got) != len(tt.allowed) {
				t.Fatalf("AllowedResultViews() = %v, want %v", got, tt.allowed)
			}
			for index, view := range tt.allowed {
				if got[index] != view {
					t.Fatalf("AllowedResultViews() = %v, want %v", got, tt.allowed)
				}
			}
			if got := tt.board.ResultView(); got != tt.defaultView {
				t.Fatalf("ResultView() = %q, want %q", got, tt.defaultView)
			}
		})
	}
}
