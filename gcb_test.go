package gcb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	res, _, _ := Go(context.Background(), "test1", func(context.Context) (int, error) {
		return 1, nil
	}, nil)

	if res != 1 {
		t.Errorf(`Breaker("test1", ...) = %v; want %v`, res, 1)
	}
}

func TestTimeout(t *testing.T) {
	res, ok, err := Go(context.Background(), t.Name(), func(context.Context) (int, error) {
		time.Sleep(5 * time.Second)
		return 1, nil
	}, nil)

	assert.Equal(t, 0, res)
	assert.False(t, ok)
	assert.Error(t, Error{Name: t.Name(), Type: ErrorTimeout, Message: "timeout after 1s"}, err)
}

func TestNoTimeout(t *testing.T) {
	res, ok, err := Go(context.Background(), t.Name(), func(context.Context) (int, error) {
		return 1, nil
	}, nil)

	assert.Equal(t, 1, res)
	assert.True(t, ok)
	assert.NoError(t, err)
}
