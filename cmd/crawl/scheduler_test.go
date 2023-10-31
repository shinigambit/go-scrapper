package crawl

import (
	"context"
	"fmt"
	"testing"
)

type mockExtractor struct {
	links  map[string][]string
	errors map[string]error
}

func (e *mockExtractor) Extract(url string) ([]string, error) {
	return e.links[url], e.errors[url]
}

func (e *mockExtractor) EffectiveDomain() string {
	return "domain.com"
}

func TestScheduler_ScheduleCrawl(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	cancelledContext, cancel := context.WithCancel(ctx)
	cancel()
	type fields struct {
		linksMap          map[string][]string
		pageErrors        map[string]error
		poolSize          int
		channelBufferSize int
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		mustNotVisit map[string]struct{}
		mustVisit    map[string]struct{}
		wantErr      string
	}{
		{
			name: "canceled context stops before visiting all urls",
			args: args{
				ctx: cancelledContext,
			},
			fields: fields{
				poolSize:          1,
				channelBufferSize: 1,
				linksMap: map[string][]string{
					"domain.com": {
						"domain.com/contact",
						"domain.com/about-us",
					},
				},
			},
			mustNotVisit: map[string]struct{}{
				"domain.com/contact":  {},
				"domain.com/about-us": {},
			},
		},
		{
			name: "should visit all urls",
			args: args{
				ctx: context.TODO(),
			},
			fields: fields{
				poolSize:          1,
				channelBufferSize: 1,
				linksMap: map[string][]string{
					"domain.com": {
						"domain.com/contact",
						"domain.com/about-us",
					},
				},
			},
			mustVisit: map[string]struct{}{
				"domain.com":          {},
				"domain.com/contact":  {},
				"domain.com/about-us": {},
			},
		},
		{
			name: "should validate poolSize greater than 0",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: "at least one worker is required to execute",
		},
		{
			name: "should validate buffer size greater than 0",
			args: args{
				ctx: context.TODO(),
			},
			fields: fields{
				poolSize: 1,
			},
			wantErr: "a buffer size of at least one is required to execute",
		},
		{
			name: "test errors",
			args: args{
				ctx: context.TODO(),
			},
			fields: fields{
				linksMap: map[string][]string{
					"domain.com": {"domain.com/contact", "domain.com/about-us"},
				},
				pageErrors: map[string]error{
					"domain.com/contact":  fmt.Errorf("first err"),
					"domain.com/about-us": fmt.Errorf("second err"),
				},
				poolSize:          1,
				channelBufferSize: 1,
			},
			mustVisit: map[string]struct{}{
				"domain.com": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				extractor:         &mockExtractor{links: tt.fields.linksMap, errors: tt.fields.pageErrors},
				poolSize:          tt.fields.poolSize,
				channelBufferSize: tt.fields.channelBufferSize,
			}
			channel, err := s.ScheduleCrawl(tt.args.ctx)
			if tt.wantErr != "" && err.Error() != tt.wantErr {
				t.Errorf("expected error (%v) but got (%v)", tt.wantErr, err)
				return
			}
			urls := map[string]struct{}{}
			errors := []error{}
			for msg := range channel {
				if msg.Url != "" {
					if _, found := tt.mustNotVisit[msg.Url]; found {
						t.Errorf("unexpected visit to url [%v]", msg.Url)
					}
					urls[msg.Url] = struct{}{}
				}
				if msg.Err != nil {
					errors = append(errors, msg.Err)
				}
			}
			if len(errors) != len(tt.fields.pageErrors) {
				t.Errorf("expected %v error(s) but found %v", len(tt.fields.pageErrors), len(errors))
			}
			for url := range tt.mustVisit {
				if _, found := urls[url]; !found {
					t.Errorf("url not visited [%v]", url)
				}
			}
		})
	}
}
