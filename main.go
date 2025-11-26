package main

import (
	"context"
	"embed"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/okunix/quietHN/hackernews"
	"github.com/okunix/quietHN/middleware"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	anyTemplate = func(name string) *template.Template {
		return template.Must(template.ParseFS(templateFS, "templates/layout.html", name))
	}
	indexTemplate = anyTemplate("templates/index.html")
)

//go:embed static/*
var staticFS embed.FS

type Validatable interface {
	Validate(ctx context.Context) (problems map[string]string)
}

func GetenvWithDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func main() {
	serverPort := GetenvWithDefault("HN_SERVER_PORT", "80")
	serverHost := GetenvWithDefault("HN_SERVER_HOST", "0.0.0.0")

	hn := hackernews.NewHackerNewsClient("https://hacker-news.firebaseio.com")
	hnCache := hackernews.NewHackerNewsClientWithCache(hn)

	router := NewRouter(hnCache)

	slog.Info("server is running", "host", serverHost, "port", serverPort)
	http.ListenAndServe(net.JoinHostPort(serverHost, serverPort), router)
}

func NewRouter(hackerNewsClient hackernews.HackerNewsClient) http.Handler {
	router := http.NewServeMux()

	router.Handle("/static/", http.FileServerFS(staticFS))
	router.Handle("GET /{$}", IndexHandler(hackerNewsClient))
	handler := middleware.Logger()(router)
	return handler
}

func IndexHandler(hnClient hackernews.HackerNewsClient) http.HandlerFunc {
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
		indexTemplate.Execute(w, templateData{stories})
	}
}
