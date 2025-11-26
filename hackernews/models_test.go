package hackernews

import "testing"

func TestExtractDomain(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{
			url:  "https://google.com/?q=test",
			want: "google.com",
		},
		{
			url:  "https://google.com?q=test",
			want: "google.com",
		},
		{
			url:  "http://example.org/path/to/page",
			want: "example.org",
		},
		{
			url:  "https://sub.domain.co.uk/some/path?param=1",
			want: "sub.domain.co.uk",
		},
		{
			url:  "https://localhost:8080/api/v1/resource",
			want: "localhost",
		},
		{
			url:  "https://192.168.1.1/dashboard",
			want: "192.168.1.1",
		},
		{
			url:  "http://my-site123.net/home?user=abc",
			want: "my-site123.net",
		},
		{
			url:  "https://www.github.com/user/repo",
			want: "github.com",
		},
		{
			url:  "https://user:pass@private.example.com:8443/path",
			want: "private.example.com",
		},
		{
			url:  "https://example.com:443/path/to/resource",
			want: "example.com",
		},
		{
			url:  "https://xn--fsq.com/",
			want: "xn--fsq.com",
		},
	}
	for _, i := range cases {
		t.Run(i.want, func(t *testing.T) {
			domain, err := extractDomain(i.url)
			if err != nil {
				t.Error(err.Error())
				return
			}
			if domain != i.want {
				t.Errorf("Domain doesn't match wanted \"%s\", got \"%s\"", i.want, domain)
				return
			}
		})
	}

}
