// Package bettererrors provides errors that are initialised with a fmt-compatible format string, so
// that values can be substituted in later.
package bettererrors

import (
    "errors"
    "fmt"
)

// Type ErrorTemplate represents an error template.
type ErrorTemplate string

// Function New creates and returns a new ErrorTemplate.
func New(format string) (et ErrorTemplate) {
    return ErrorTemplate(format)
}

// Function Format runs fmt.Sprintf using the format string specified during creation and the specified
// arguments, then returns the created error.
func (et ErrorTemplate) Format(args ...interface{}) (err error) {
    return errors.New(fmt.Sprintf(string(et), args...))
}
