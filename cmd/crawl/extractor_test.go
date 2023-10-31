package crawl

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

type mockGetter struct {
	s string
}

func (g *mockGetter) Get(url string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(g.s)), nil
}

func TestClient_Request(t *testing.T) {
	t.Parallel()
	type fields struct {
		domain   string
		response string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantLinks []string
		wantErr   bool
	}{
		{
			name: "convert all links to full urls",
			args: args{
				url: "https://site.to.crawl.com/x/",
			},
			fields: fields{
				domain: "https://site.to.crawl.com",
				response: `
						<html>
							<body>
								<a href="/absolute-link.html#a"/>
								<a href='relative-link.xml/#a'>some text</a>
							</body>
						</html>
						`,
			},
			wantLinks: []string{"https://site.to.crawl.com/absolute-link.html", "https://site.to.crawl.com/x/relative-link.xml"},
		},
		{
			name: "skip external domain links",
			args: args{
				url: "https://site.to.crawl.com",
			},
			fields: fields{
				domain: "https://site.to.crawl.com/",
				response: `
						<html>
							<body>
								<a href="https://site.to.crawl.com/first.html"/>
								<a href='http://site.to.crawl.com/different-protocol.html'></a>
								<a href="https://same.protocol.different.domain.com/"/>
							</body>
						</html>`,
			},
			wantLinks: []string{"https://site.to.crawl.com/first.html"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := NewUrlLinkExtractor(tt.fields.domain, &mockGetter{s: tt.fields.response})
			gotLinks, err := c.Extract(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLinks, tt.wantLinks) {
				t.Errorf("Client.Request() = %v, want %v", gotLinks, tt.wantLinks)
			}
		})
	}
}
