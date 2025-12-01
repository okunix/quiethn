package hackernews

import (
	"context"
	"errors"
	"regexp"
	"time"
)

type NewsItem struct {
	Id     uint32    `json:"id"`
	Title  string    `json:"title"`
	URL    string    `json:"url"`
	Domain string    `json:"domain"`
	Type   string    `json:"type"`
	Time   time.Time `json:"time"`
}

var (
	domainRegex          = regexp.MustCompile(`^https?://(?:.+?@)?(?:www.)?([^\s/?:]+)(?:[/?:]|$)`)
	trimSpacesRegex      = regexp.MustCompile(`^\s+|\s+$`)
	htmlTagRegex         = regexp.MustCompile(`<[^>]*>`)
	sequentialWhitespace = regexp.MustCompile(`\s+`)
)

func extractDomain(url string) (string, error) {
	regexSubmatchList := domainRegex.FindStringSubmatch(url)
	if len(regexSubmatchList) < 2 {
		return "", errors.New("no domain detected")
	}
	return regexSubmatchList[1], nil
}

func extractTextFromHTML(html string) string {
	s := htmlTagRegex.ReplaceAllString(html, "")
	s = sequentialWhitespace.ReplaceAllString(s, " ")
	s = trimSpacesRegex.ReplaceAllString(s, "")
	return s
}

func NewNewsItem(id uint32, title, url, itemType string, time time.Time) (NewsItem, error) {
	domain, err := extractDomain(url)
	if err != nil {
		return NewsItem{}, err
	}
	newsItem := NewsItem{
		Id:     id,
		Title:  extractTextFromHTML(title),
		URL:    url,
		Domain: domain,
		Type:   itemType,
		Time:   time,
	}
	return newsItem, nil
}

func (ni *NewsItem) Validate(ctx context.Context) (problems map[string]string) {
	panic("unimplemented")
}

type NewsItemResponse struct {
	Id      uint32 `json:"id"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	Type    string `json:"type"`
	Time    int64  `json:"time"`
	Deleted bool   `json:"deleted"`
}
