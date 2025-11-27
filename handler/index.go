package handler

import (
	"net/http"

	"github.com/okunix/quietHN/hackernews"
	"github.com/okunix/quietHN/templates"
)

func Index(hnClient hackernews.HackerNewsClient) http.HandlerFunc {
	type templateData struct {
		Stories []*hackernews.NewsItem
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		stories, err := hnClient.GetTopStories(ctx, 30)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templates.Index.Execute(w, templateData{stories})
	}
}
