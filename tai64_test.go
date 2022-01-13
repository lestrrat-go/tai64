package tai64_test

import (
	"testing"
	"time"

	"github.com/lestrrat-go/tai64"
	"github.com/stretchr/testify/assert"
)

func TestParseNLabel(t *testing.T) {
	testcases := []struct {
		ExpectedTime  time.Time
		ExpectedBytes []byte
		Input         []byte
		Error         bool
	}{
		{
			Input:        []byte(`4000000037c219bf2ef02e94`),
			ExpectedTime: time.Unix(935467455, 787492500),
		},
		{
			Input:         []byte(`@4000000037c219bf2ef02e94`),
			ExpectedBytes: []byte(`4000000037c219bf2ef02e94`),
			ExpectedTime:  time.Unix(935467455, 787492500),
		},
		{
			Input: []byte(`!4000000037c219bf2ef02e94`),
			Error: true,
		},
		{
			Input: []byte(`@4000000037c219bf2ef02e`),
			Error: true,
		},
		{
			Input: []byte(`@4000000037c219bf2ef02e9494`),
			Error: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(string(tc.Input), func(t *testing.T) {
			n, err := tai64.ParseNLabel(tc.Input)
			if tc.Error {
				if !assert.Error(t, err, `tai64.ParseNLabel should fail`) {
					return
				}
			} else {
				if !assert.NoError(t, err, `tai64.ParseNLabel should succeed`) {
					return
				}

				tv := n.Time()
				if !assert.Equal(t, tc.ExpectedTime, tv) {
					return
				}

				var dst [tai64.NLabelSize]byte
				n.Format(dst[:])
				expected := tc.ExpectedBytes
				if len(expected) == 0 {
					expected = tc.Input
				}
				if !assert.Equal(t, expected, dst[:], `serialized form of tai64.N should match`) {
					return
				}

				txt, err := n.MarshalText()
				if !assert.NoError(t, err, `n.MarshalText should succeed`) {
					return
				}
				if !assert.Equal(t, expected, txt, `serialized form of tai64.N should match (via MarshalText())`) {
					return
				}
			}
		})
	}
}
