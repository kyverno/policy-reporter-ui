package api

import (
	"fmt"
	"sort"

	"github.com/kyverno/policy-reporter-ui/pkg/core/client"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

func MapSourceCategoryTreeToNavi(sources []client.SourceCategoryTree) []NavigationItem {
	sourceBoards := make([]NavigationItem, 0)
	if len(sources) == 1 {
		for _, category := range sources[0].Categories {
			sourceBoards = append(sourceBoards, NavigationItem{
				Title:    utils.Title(sources[0].Name),
				Subtitle: category.Name,
				Path:     fmt.Sprintf("/source/%s/%s", sources[0].Name, category.Name),
			})
		}

		return sourceBoards
	}

	for _, source := range sources {
		count := len(source.Categories)
		if count == 0 {
			sourceBoards = append(sourceBoards, NavigationItem{
				Title: utils.Title(source.Name),
				Path:  fmt.Sprintf("/source/%s", source.Name),
			})
		} else if count == 1 {
			sourceBoards = append(sourceBoards, NavigationItem{
				Title:    utils.Title(source.Name),
				Subtitle: source.Categories[0].Name,
				Path:     fmt.Sprintf("/source/%s/%s", source.Name, source.Categories[0].Name),
			})
		} else {
			item := NavigationItem{
				Title:    utils.Title(source.Name),
				Path:     fmt.Sprintf("/source/%s", source.Name),
				Children: make([]NavigationItem, 0, len(source.Categories)),
			}

			for _, category := range source.Categories {
				item.Children = append(item.Children, NavigationItem{
					Title: category.Name,
					Path:  fmt.Sprintf("/source/%s/%s", source.Name, category.Name),
				})
			}

			sourceBoards = append(sourceBoards, item)
		}
	}

	sort.SliceStable(sourceBoards, func(a, b int) bool {
		return sourceBoards[a].Title < sourceBoards[b].Title
	})

	return sourceBoards
}

func MapCustomBoardsToNavi(boards map[string]CustomBoard) []NavigationItem {
	customBoards := make([]NavigationItem, 0, len(boards))
	for _, board := range boards {
		customBoards = append(customBoards, NavigationItem{
			Title: board.Name,
			Path:  fmt.Sprintf("/custom-boards/%s", board.ID),
		})
	}

	sort.SliceStable(customBoards, func(a, b int) bool {
		return customBoards[a].Title < customBoards[b].Title
	})

	return customBoards
}
