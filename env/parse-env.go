package env

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
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
			return errors.Wrapf(err, "failed to parse %q", key)
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
		return errors.Wrapf(err, "unable to parse %q as int", s)
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
		return errors.Wrapf(err, "unable to parse %q as float", s)
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
		return errors.Wrapf(err, "unable to parse %q as bool", s)
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

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Environment -s _mock.go

type Environment interface {
	Getenv(string) string
}

type osEnv struct {
	// defer to the os package
}

func (e *osEnv) Getenv(name string) string {
	return os.Getenv(name)
}

var OS Environment = &osEnv{}
