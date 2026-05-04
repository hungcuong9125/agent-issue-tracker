package ait

import (
	"errors"
	"fmt"
)

// exitError signals a specific shell exit code from a command handler.
// When err is nil the dispatcher skips the JSON error envelope — used for
// clean non-zero exits like 'self-update --check' reporting that a newer
// version is available without that being an error condition.
type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return fmt.Sprintf("exit %d", e.code)
}

func (e *exitError) Unwrap() error { return e.err }

// ExitWithCode wraps err so the binary exits with the given code. A nil err
// produces a silent non-zero exit (no envelope written).
func ExitWithCode(code int, err error) error {
	return &exitError{code: code, err: err}
}

// ExitCode reports whether err carries a specific exit code. main uses this
// to translate command results into shell exit status.
func ExitCode(err error) (int, bool) {
	var e *exitError
	if errors.As(err, &e) {
		return e.code, true
	}
	return 0, false
}

// SilentExit reports whether err is an exitError whose embedded cause is
// nil — main uses this to suppress the JSON envelope for clean signal-only
// exits.
func SilentExit(err error) bool {
	var e *exitError
	if errors.As(err, &e) {
		return e.err == nil
	}
	return false
}
