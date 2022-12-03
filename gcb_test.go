package gcb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	t.Parallel()

	res, prim, _ := Go(context.Background(), "test1", func(context.Context) (int, error) {
		return 1, nil
	}, nil)

	assert.True(t, prim)

	if res != 1 {
		t.Errorf(`Go(ctx, "test1", ...) = %v; want %v`, res, 1)
	}
}

func TestFallback(t *testing.T) {
	t.Parallel()

	res, prim, _ := Go(context.Background(), "test1", func(context.Context) (int, error) {
		return 1, fmt.Errorf("This is an error")
	}, func(context.Context) (int, error) {
		return 2, nil
	})

	assert.True(t, prim)

	if res != 2 {
		t.Errorf(`Go(ctx, "test1", ...) = %v; want %v`, res, 1)
	}
}

func TestTimeout(t *testing.T) {
	t.Parallel()
	res, ok, err := Go(context.Background(), t.Name(), func(context.Context) (int, error) {
		time.Sleep(5 * time.Second)
		return 1, nil
	}, nil)

	assert.Equal(t, 0, res)
	assert.False(t, ok)
	assert.Error(t, Error{Name: t.Name(), Type: ErrorTimeout, Message: "timeout after 1s"}, err)
}

func TestNoTimeout(t *testing.T) {
	t.Parallel()
	res, ok, err := Go(context.Background(), t.Name(), func(context.Context) (int, error) {
		return 1, nil
	}, nil)

	assert.Equal(t, 1, res)
	assert.True(t, ok)
	assert.NoError(t, err)
}
