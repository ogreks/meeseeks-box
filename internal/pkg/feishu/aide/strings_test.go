package aide

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntherCutPrefix(t *testing.T) {
	testCase := []struct {
		name string
		args struct {
			s       string
			prefixs []string
		}
		want   string
		result bool
	}{
		{
			name: "prefix match",
			args: struct {
				s       string
				prefixs []string
			}{
				s:       "检查 1",
				prefixs: []string{"检查"},
			},
			want:   "1",
			result: true,
		},
		{
			name: "prefix match command",
			args: struct {
				s       string
				prefixs []string
			}{
				s:       "/check 1",
				prefixs: []string{"/check "},
			},
			want:   "1",
			result: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			want, result := EitherCutPrefix(tc.args.s, tc.args.prefixs...)
			assert.Equal(t, tc.want, want)
			assert.Equal(t, tc.result, result)
		})
	}
}

func TestEitherTrimEqual(t *testing.T) {
	testCase := []struct {
		name string
		args struct {
			s       string
			prefixs []string
		}
		want   string
		result bool
	}{
		{
			name: "prefix match command",
			args: struct {
				s       string
				prefixs []string
			}{
				s:       "/version",
				prefixs: []string{"/version"},
			},
			want:   "",
			result: true,
		},
		{
			name: "prefix match",
			args: struct {
				s       string
				prefixs []string
			}{
				s:       " version ",
				prefixs: []string{"/version"},
			},
			want:   " version ",
			result: false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			want, result := EitherTrimEqual(tc.args.s, tc.args.prefixs...)
			assert.Equal(t, tc.want, want)
			assert.Equal(t, tc.result, result)
		})
	}
}
