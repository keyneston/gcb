package gcb

import (
	"context"
	"fmt"
	"time"
)

type Error struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%#q: %v", e.Name, e.Message)
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
	circuit := GetCircuit(name)

	return processCircuit(
		ctx,
		name,
		settings,
		circuit,
		primary,
		fallback,
	)
}

func processCircuit[T any](
	ctx context.Context,
	name string,
	settings Settings,
	circuit *Circuit,
	primary func(context.Context) (T, error),
	fallback func(context.Context) (T, error),
) (T, bool, error) {
	res, err := handleFunction(ctx, name, settings, circuit, primary)
	// If fallback is not defined or we received no errors return whatever we got
	if fallback == nil || err == nil {
		return res, true, err
	}

	res, err = handleFunction(ctx, name, settings, circuit, fallback)
	return res, false, err
}

func handleFunction[T any](
	ctx context.Context,
	name string,
	settings Settings,
	circuit *Circuit,
	fn func(context.Context) (T, error),
) (T, error) {
	timer := time.NewTimer(settings.Timeout.D())
	defer cleanTimer(timer)

	tChan := make(chan T, 1)
	errChan := make(chan error, 1)

	// TODO thing here to setup context

	go func() {
		res, err := fn(ctx)
		tChan <- res
		errChan <- err

		defer close(tChan)
		defer close(errChan)
	}()

	var ret T
	var retErr error

	var tReceived, errReceived bool

	for {
		select {
		case t, ok := <-tChan:
			// If this is a channel closed message don't update the return value.
			if ok {
				ret = t
			}
			tReceived = true
		case err, ok := <-errChan:
			// If this is a channel closed message don't update the return value.
			if ok {
				retErr = err
			}
			errReceived = true
		case <-timer.C:
			return ret, Error{
				Name:    name,
				Type:    ErrorTimeout,
				Message: fmt.Sprintf("timeout after %v", settings.Timeout),
			}
		case <-ctx.Done():
			return ret, Error{
				Name:    name,
				Type:    ErrorContext,
				Message: fmt.Sprintf("context done"),
			}
		}

		if tReceived && errReceived {
			return ret, retErr
		}
	}

	// TODO: handle fallback function
	return ret, retErr
}
