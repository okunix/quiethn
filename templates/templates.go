package templates

import (
	"embed"
	"text/template"
)

//go:embed *.html
var templateFS embed.FS

var funcMap = template.FuncMap{
	"inc": func(n int) int {
		return n + 1
	},
	"trunc": func(n int, add string, s string) string {
		if len(s) <= n+len(add) {
			return s
		}
		return s[:n] + add
	},
}

func anyTemplate(templatePath string) *template.Template {
	return template.Must(
		template.New("layout.html").Funcs(funcMap).
			ParseFS(templateFS, templatePath, "layout.html"),
	)
}

var (
	Index = anyTemplate("index.html")
)
