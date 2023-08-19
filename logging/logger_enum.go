// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.7
// Revision: bf63e108589bbd2327b13ec2c5da532aad234029
// Build Date: 2023-07-25T23:27:55Z
// Built By: goreleaser

package logging

import (
	"errors"
	"fmt"
)

const (
	// EnvDevelopment is a Env of type Development.
	EnvDevelopment Env = "development"
	// EnvProduction is a Env of type Production.
	EnvProduction Env = "production"
)

var ErrInvalidEnv = errors.New("not a valid Env")

// String implements the Stringer interface.
func (x Env) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Env) IsValid() bool {
	_, err := ParseEnv(string(x))
	return err == nil
}

var _EnvValue = map[string]Env{
	"development": EnvDevelopment,
	"production":  EnvProduction,
}

// ParseEnv attempts to convert a string to a Env.
func ParseEnv(name string) (Env, error) {
	if x, ok := _EnvValue[name]; ok {
		return x, nil
	}
	return Env(""), fmt.Errorf("%s is %w", name, ErrInvalidEnv)
}
