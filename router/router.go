package router

import (
	"io/fs"
	"net/http"

	"github.com/okunix/quiethn/hackernews"
	"github.com/okunix/quiethn/handler"
	"github.com/okunix/quiethn/middleware"
)

func NewRouter(
	hackerNewsRepo hackernews.HackerNewsRepo,
	staticFS fs.FS,
) http.Handler {
	router := http.NewServeMux()

	router.Handle("/static/", http.FileServerFS(staticFS))
	router.Handle("GET /{$}", handler.Index(hackerNewsRepo))

	handler := middleware.Logger()(router)
	handler = middleware.RealIP()(handler)
	//handler = middleware.NoCache()(handler)
	return handler
}
