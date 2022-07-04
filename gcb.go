package gcb

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

var defaultSettings *atomic.Value = &atomic.Value{}

func init() {
	SetDefaultSettings(Settings{
		Timeout:         time.Second,
		FallbackTimeout: time.Second,
	})
}

func GetDefaultSettings() Settings {
	return *defaultSettings.Load().(*Settings)
}

func SetDefaultSettings(s Settings) {
	defaultSettings.Store(&s)
}

type Error struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%#q: %v", e.Name, e.Message)
}

type Settings struct {
	Timeout         time.Duration
	FallbackTimeout time.Duration

	Percent float32
}

// Go calls your primary function, and if that fails calls your fallback
// function if provided.
//
// Returns the output of the primary or fallback, along with true if the result
// was from the primary function, or false if it was from the fallback function
// or the library itself.
func Go[T any](
	ctx context.Context,
	name string,
	primary func(context.Context) (T, error),
	fallback func(context.Context) (T, error),
) (T, bool, error) {
	settings := GetDefaultSettings()
	timer := time.NewTimer(settings.Timeout)
	defer cleanTimer(timer)

	tChan := make(chan T, 1)
	errChan := make(chan error, 1)
	defer close(tChan)
	defer close(errChan)

	// TODO thing here to setup context

	go func() {
		res, err := primary(ctx)
		tChan <- res
		errChan <- err
	}()

	var t T
	var err error

	var tReceived, errReceived bool

	for {
		select {
		case t = <-tChan:
			tReceived = true
			if errReceived {
				return t, true, err
			}
		case err = <-errChan:
			errReceived = true
			if tReceived {
				return t, true, err
			}
		case <-timer.C:
			return t, false, Error{
				Name:    name,
				Type:    ErrorTimeout,
				Message: fmt.Sprintf("timeout after %v", settings.Timeout),
			}
		case <-ctx.Done():
			return t, false, Error{
				Name:    name,
				Type:    ErrorContext,
				Message: fmt.Sprintf("context done"),
			}
		}
	}

	// TODO: handle fallback function
	return t, true, err
}
