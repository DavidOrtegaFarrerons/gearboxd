package data

import (
	"gearboxd/internal/assert"
	"gearboxd/internal/validator"
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

func TestValidateFilters(t *testing.T) {
	testCases := []struct {
		name           string
		page           int
		pageSize       int
		sort           string
		sortSafelist   []string
		valid          bool
		expectedErrors map[string]string
	}{
		{
			name:           "Correct data with no errors",
			page:           1,
			pageSize:       10,
			sort:           "make",
			sortSafelist:   []string{"make"},
			valid:          true,
			expectedErrors: nil,
		},
		{
			name:           "Page is less than 0",
			page:           -1,
			pageSize:       10,
			sort:           "make",
			sortSafelist:   []string{"make"},
			valid:          false,
			expectedErrors: map[string]string{"page": "must be greater than 0"},
		},
		{
			name:           "Page is higher than 1.000.000",
			page:           1_000_001,
			pageSize:       10,
			sort:           "make",
			sortSafelist:   []string{"make"},
			valid:          false,
			expectedErrors: map[string]string{"page": "must be a maximum of 1 million"},
		},
		{
			name:           "PageSize is less than 0",
			page:           1,
			pageSize:       -1,
			sort:           "make",
			sortSafelist:   []string{"make"},
			valid:          false,
			expectedErrors: map[string]string{"page_size": "must be greater than 0"},
		},
		{
			name:           "PageSize is higher than 1.000.000",
			page:           10,
			pageSize:       1_000_001,
			sort:           "make",
			sortSafelist:   []string{"make"},
			valid:          false,
			expectedErrors: map[string]string{"page_size": "must be a maximum of 1 million"},
		},
		{
			name:           "Sort is inside sort safelist",
			page:           2,
			pageSize:       10,
			sort:           "make",
			sortSafelist:   []string{"make"},
			valid:          true,
			expectedErrors: nil,
		},
		{
			name:           "Sort is not inside sort safelist",
			page:           2,
			pageSize:       10,
			sort:           "make",
			sortSafelist:   []string{"model"},
			valid:          false,
			expectedErrors: map[string]string{"sort": "invalid sort value"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f := Filters{
				Page:         tt.page,
				PageSize:     tt.pageSize,
				Sort:         tt.sort,
				SortSafelist: tt.sortSafelist,
			}

			v := validator.New()

			ValidateFilters(v, f)

			assert.Equal(t, v.Valid(), tt.valid)
			if !tt.valid {
				assert.Equal(t, v.Errors, tt.expectedErrors)
			}
		})
	}
}
