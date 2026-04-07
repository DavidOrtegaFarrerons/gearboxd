package data

import (
	"gearboxd/internal/assert"
	"testing"
)

func TestSortColumn(t *testing.T) {
	testCases := []struct {
		name     string
		safeList []string
		sort     string
		want     string
	}{
		{
			name:     "Sort without dash is inside safeList",
			safeList: []string{"make", "model", "horsepower"},
			sort:     "make",
			want:     "make",
		},
		{
			name:     "Sort with dash is inside safeList, sort without dash is returned",
			safeList: []string{"-make", "model", "horsepower"},
			sort:     "-make",
			want:     "make",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f := Filters{
				SortSafelist: tt.safeList,
				Sort:         tt.sort,
			}

			got := f.sortColumn()

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestSortDirection(t *testing.T) {
	testCases := []struct {
		name string
		sort string
		want string
	}{
		{
			name: "Sort does not contain dash",
			sort: "make",
			want: "ASC",
		},
		{
			name: "Sort contains a dash",
			sort: "-make",
			want: "DESC",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f := Filters{
				Sort: tt.sort,
			}

			got := f.sortDirection()

			assert.Equal(t, got, tt.want)
		})
	}
}
