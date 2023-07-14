package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/shoenig/go-conceal"
)

// A Variable represents a set environment variable
type Variable string

func (v Variable) String() string {
	return "{" + string(v) + "}"
}

func (v Variable) Name() string {
	return string(v)
}

type Schema map[Variable]Parser

func ParseOS(schema Schema) error {
	return Parse(OS, schema)
}

func Parse(environment Environment, schema Schema) error {
	for key, parser := range schema {
		value := environment.Getenv(key.Name())
		if err := parser.Parse(value); err != nil {
			return fmt.Errorf("failed to parse %q: %w", key, err)
		}
	}
	return nil
}

type Parser interface {
	Parse(string) error
}

type stringParser struct {
	required    bool
	destination *string
}

func (sp *stringParser) Parse(s string) error {
	if sp.required && s == "" {
		return errors.New("missing")
	} else if s == "" {
		return nil
	}

	*sp.destination = s
	return nil
}

func String(s *string, required bool) Parser {
	return &stringParser{
		required:    required,
		destination: s,
	}
}

type secretParser struct {
	required    bool
	destination **conceal.Text
}

func (sp *secretParser) Parse(s string) error {
	if sp.required && s == "" {
		return errors.New("missing")
	} else if s == "" {
		return nil
	}

	text := conceal.New(s)
	*sp.destination = text
	return nil
}

func Secret(s **conceal.Text, required bool) Parser {
	return &secretParser{
		required:    required,
		destination: s,
	}
}

type intParser struct {
	required    bool
	destination *int
}

func (ip *intParser) Parse(s string) error {
	if ip.required && s == "" {
		return errors.New("missing")
	} else if s == "" {
		return nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("unable to parse %q as int: %w", s, err)
	}
	*ip.destination = i
	return nil
}

func Int(i *int, required bool) Parser {
	return &intParser{
		required:    required,
		destination: i,
	}
}

type floatParser struct {
	required    bool
	destination *float64
}

func (fp *floatParser) Parse(s string) error {
	if fp.required && s == "" {
		return errors.New("missing")
	} else if s == "" {
		return nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("unable to parse %q as float: %w", s, err)
	}
	*fp.destination = f
	return nil
}

func Float(f *float64, required bool) Parser {
	return &floatParser{
		required:    required,
		destination: f,
	}
}

type boolParser struct {
	required    bool
	destination *bool
}

func (bp *boolParser) Parse(s string) error {
	if bp.required && s == "" {
		return errors.New("missing")
	} else if s == "" {
		return nil
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("unable to parse %q as bool: %w", s, err)
	}
	*bp.destination = b
	return nil
}

func Bool(b *bool, required bool) Parser {
	return &boolParser{
		required:    required,
		destination: b,
	}
}

// Environment is something that implements Getenv().
//
// Most use cases can simply make use of the OS implementation which is backed
// by the standard library os package.
type Environment interface {
	Getenv(string) string
}

type osEnv struct {
	// defer to the os package
}

func (e *osEnv) Getenv(name string) string {
	return os.Getenv(name)
}

// OS is an implementation of Environment that uses the standard library os
// package to retrieve actual environment variables.
var OS Environment = &osEnv{}
