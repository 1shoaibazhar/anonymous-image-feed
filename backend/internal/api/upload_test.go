package api

import (
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "empty string",
			raw:  "",
			want: []string{},
		},
		{
			name: "single tag",
			raw:  "cats",
			want: []string{"cats"},
		},
		{
			name: "multiple tags",
			raw:  "cats,dogs,birds",
			want: []string{"cats", "dogs", "birds"},
		},
		{
			name: "trims whitespace around tags",
			raw:  " cats , dogs ,birds ",
			want: []string{"cats", "dogs", "birds"},
		},
		{
			name: "drops empty entries from stray commas",
			raw:  "cats,,dogs,",
			want: []string{"cats", "dogs"},
		},
		{
			name: "whitespace only input yields no tags",
			raw:  "   ",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTags(tt.raw)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags(%q) = %#v, want %#v", tt.raw, got, tt.want)
			}
		})
	}
}
