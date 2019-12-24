package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeVariables(t *testing.T) {
	type testcase struct {
		envs map[string]string
		want map[string]string
	}

	testcases := []testcase{
		{
			envs: map[string]string{
				"FOO": "cloud-secrets://nop/foo",
				"BAR": "cloud-secrets://nop/bar",
				"BAZ": "baz",
			},
			want: map[string]string{
				"FOO": "cloud-secrets://decoded/foo",
				"BAR": "cloud-secrets://decoded/bar",
			},
		},
	}

	for _, tc := range testcases {
		d := NewDecoder()
		got, err := d.DecodeVariables(tc.envs)

		if !assert.NoError(t, err) {
			continue
		}

		assert.Equal(t, tc.want, got)
	}
}
