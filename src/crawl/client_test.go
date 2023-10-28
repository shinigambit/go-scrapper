package crawl

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestClient_Request(t *testing.T) {
	t.Parallel()
	type fields struct {
		domain string
		getUrl func(string) (*http.Response, error)
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
			name: "",
			args: args{
				url: "https://site.to.scrap.com/x/",
			},
			fields: fields{
				domain: "https://site.to.scrap.com/",
				getUrl: func(s string) (*http.Response, error) {
					return &http.Response{
						Body: io.NopCloser(strings.NewReader(`
						<html>
							<body>
								<a href="/link.html#a"/>
								<a href="link.xml/#a"></a>
							</body>
						</html>
						`)),
					}, nil
				},
			},
			wantLinks: []string{"https://site.to.scrap.com/link.html", "https://site.to.scrap.com/x/link.xml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := NewClient(tt.fields.domain)
			c.getUrl = tt.fields.getUrl
			gotLinks, err := c.Request(tt.args.url)
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
