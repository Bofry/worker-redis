package internal

import "fmt"

var _ error = new(FatalError)

type FatalError struct {
	err interface{}
}

// Error implements error.
func (e *FatalError) Error() string {
	switch err := e.err.(type) {
	case error:
		return err.Error()
	case string:
		return err
	case fmt.Stringer:
		return err.String()
	}
	return fmt.Sprintf("%+v", e.err)
}
