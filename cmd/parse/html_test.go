package parse

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestPageLinks(t *testing.T) {
	t.Parallel()
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name      string
		args      args
		wantLinks []string
		wantErr   bool
	}{
		{
			name: "parses both link tags",
			args: args{
				reader: strings.NewReader(`
				<html>
					<body>
						<a href="/go.html#">link name</a>
						<a href='go.html'/>
					</body>
				</html>
				`),
			},
			wantLinks: []string{"/go.html#", "go.html"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotLinks, err := PageLinks(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLinks, tt.wantLinks) {
				t.Errorf("PageLinks() = %v, want %v", gotLinks, tt.wantLinks)
			}
		})
	}
}
