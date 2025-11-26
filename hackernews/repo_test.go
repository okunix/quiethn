package hackernews

import "testing"

func TestGetItemById(t *testing.T) {
	ctx := t.Context()
	repo := NewHackerNewsClient("https://hacker-news.firebaseio.com")
	newsItem, err := repo.GetItemById(ctx, 46055421)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%+v\n", newsItem)
}
