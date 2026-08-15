package repository

import (
	"reflect"
	"testing"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name        string
		images      []Image
		pageSize    int
		wantLen     int
		wantHasMore bool
	}{
		{
			name:        "fewer rows than page size",
			images:      []Image{{ID: "a"}, {ID: "b"}},
			pageSize:    5,
			wantLen:     2,
			wantHasMore: false,
		},
		{
			name:        "exactly page size",
			images:      []Image{{ID: "a"}, {ID: "b"}},
			pageSize:    2,
			wantLen:     2,
			wantHasMore: false,
		},
		{
			name:        "one extra row signals more pages",
			images:      []Image{{ID: "a"}, {ID: "b"}, {ID: "c"}},
			pageSize:    2,
			wantLen:     2,
			wantHasMore: true,
		},
		{
			name:        "empty input",
			images:      []Image{},
			pageSize:    2,
			wantLen:     0,
			wantHasMore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, hasMore := paginate(tt.images, tt.pageSize)
			if len(got) != tt.wantLen {
				t.Errorf("paginate() returned %d images, want %d", len(got), tt.wantLen)
			}
			if hasMore != tt.wantHasMore {
				t.Errorf("paginate() hasMore = %v, want %v", hasMore, tt.wantHasMore)
			}
		})
	}
}

func TestAttachTags(t *testing.T) {
	tests := []struct {
		name        string
		images      []Image
		tagsByImage map[string][]string
		want        []Image
	}{
		{
			name:   "image with tags",
			images: []Image{{ID: "a"}},
			tagsByImage: map[string][]string{
				"a": {"cats", "cute"},
			},
			want: []Image{{ID: "a", Tags: []string{"cats", "cute"}}},
		},
		{
			name:        "image missing from map gets empty slice, not nil",
			images:      []Image{{ID: "b"}},
			tagsByImage: map[string][]string{},
			want:        []Image{{ID: "b", Tags: []string{}}},
		},
		{
			name:   "mixed page of tagged and untagged images",
			images: []Image{{ID: "a"}, {ID: "b"}, {ID: "c"}},
			tagsByImage: map[string][]string{
				"a": {"cats"},
				"c": {"funny", "cats"},
			},
			want: []Image{
				{ID: "a", Tags: []string{"cats"}},
				{ID: "b", Tags: []string{}},
				{ID: "c", Tags: []string{"funny", "cats"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attachTags(tt.images, tt.tagsByImage)
			if !reflect.DeepEqual(tt.images, tt.want) {
				t.Errorf("attachTags() = %#v, want %#v", tt.images, tt.want)
			}
		})
	}
}
