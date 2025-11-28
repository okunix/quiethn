package handler

import (
	"net/http"

	"github.com/okunix/quiethn/hackernews"
	"github.com/okunix/quiethn/templates"
)

func Index(hnRepo hackernews.HackerNewsRepo) http.HandlerFunc {
	type templateData struct {
		Stories []*hackernews.NewsItem
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		stories, err := hnRepo.GetTopStories(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templates.Index.Execute(w, templateData{stories})
	}
}
