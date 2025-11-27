package router

import (
	"io/fs"
	"net/http"

	"github.com/okunix/quietHN/hackernews"
	"github.com/okunix/quietHN/handler"
	"github.com/okunix/quietHN/middleware"
)

func NewRouter(
	hackerNewsClient hackernews.HackerNewsClient,
	staticFS fs.FS,
) http.Handler {
	router := http.NewServeMux()

	router.Handle("/static/", http.FileServerFS(staticFS))
	router.Handle("GET /{$}", handler.Index(hackerNewsClient))

	handler := middleware.Logger()(router)
	handler = middleware.RealIP()(handler)
	//handler = middleware.NoCache()(handler)
	return handler
}
