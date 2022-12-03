package rolling

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCount = 60

func TestRolling(t *testing.T) {
	type testCase struct {
		name     string
		fn       func(r *rolling)
		expected bucket
	}

	testCases := []testCase{
		{
			name: "success",
			fn: func(r *rolling) {
				r.Succeed()
			},
			expected: bucket{attempts: 1, failures: 0},
		},
		{
			name: "failures",
			fn: func(r *rolling) {
				r.Fail()
			},
			expected: bucket{attempts: 1, failures: 1},
		},
		{
			name: "succeed-and-fail",
			fn: func(r *rolling) {
				r.Succeed()
				r.Fail()
			},
			expected: bucket{attempts: 2, failures: 1},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			r := New(testCount)

			require.NotNil(t, r.Buckets)
			require.NotNil(t, r.Mutex)
			require.NotNil(t, r.time)
			require.Equal(t, r.count, testCount)

			r.time = func() time.Time {
				// return an empty time so that it (<x> % time.Seconds == 0)
				return time.Time{}
			}

			c.fn(r)
			assert.Equal(t, r.Buckets[0], &c.expected)
		})
	}
}
