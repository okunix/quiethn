package hackernews

import (
	"testing"
)

func TestGetItemById(t *testing.T) {
	ctx := t.Context()
	repo := NewHackerNewsClient("https://hacker-news.firebaseio.com")
	newsItem, err := repo.GetItemById(ctx, 46055421)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("%+v\n", newsItem)
}

func TestGetTopItems(t *testing.T) {
	ctx := t.Context()
	repo := NewHackerNewsClient("https://hacker-news.firebaseio.com")
	ids, err := repo.GetTopItems(ctx)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(ids) != 500 {
		t.Error("number of ids must be equal to 500")
		return
	}
	t.Logf("len(ids):%d %+v\n", len(ids), ids)
}
