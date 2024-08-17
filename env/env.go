package env

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"

	"github.com/shoenig/go-conceal"
)

// A Variable represents an environment variable.
type Variable string

func (v Variable) String() string {
	return "{" + string(v) + "}"
}

func (v Variable) Name() string {
	return string(v)
}

// Schema is used to describe how to parse a set of environment variables.
type Schema map[Variable]Parser

// ParseOS is a convenience function for parsing the given Schema of environment
// variables using the environment variables accessed by the standard libraray
// os package. If the values of environment variables do not match the schema,
// or required variables are missing, an error is returned.
func ParseOS(schema Schema) error {
	return Parse(OS, schema)
}

// ParseFile is a convenience function for parsing the given Schema of environment
// variables using the given .env file path. The contents of the file are read
// and interpreted as key=value pairs, one per line. If the environment variable
// contents of the file do not match the schema, or required variables are missing,
// an error is returned.
func ParseFile(path string, schema Schema) error {
	return Parse(File(path), schema)
}

// ParseMap is a convenience function for parsing the given Schema of environment
// variables using the given map. The contents of the map are inferred as
// key=value pairs. If the contents of the map do not match the schema, or
// required variables are missing, an error is returned.
func ParseMap(m map[string]string, schema Schema) error {
	return Parse(Map(m), schema)
}

// Parse uses the given Schema to parse the environment variables in the given
// Environment. If the values of environment variables in Environment do not
// match the schema, or required variables are missing, an error is returned.
func Parse(environment Environment, schema Schema) error {
	for key, parser := range schema {
		value := environment.Getenv(key.Name())
		if err := parser.Parse(value); err != nil {
			return fmt.Errorf("failed to parse %q: %w", key, err)
		}
	}
	return nil
}

// The Parser interface is what must be implemented to support decoding an
// environment variable into a custom type.
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

// String is used to extract an environment variable into a Go string. If
// required is true, then an error is returned if the environment variable is
// not set or is empty.
func String(s *string, required bool) Parser {
	return &stringParser{
		required:    required,
		destination: s,
	}
}

// StringOr is used to extract an environment variable into a Go string. If
// the environment variable is not set or is empty, then the alt value is used
// instead.
func StringOr(s *string, alt string) Parser {
	*s = alt
	return &stringParser{
		required:    false,
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

// Secret is used to extract an environment variable into a Go concel.Text
// object. If required is true, then an error is returned if the environment
// variable is not set or is empty.
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

// Int is used to extract an environment variable into a Go int. If required
// is true, then an error is returned if the environment variable is not set or
// is empty.
func Int(i *int, required bool) Parser {
	return &intParser{
		required:    required,
		destination: i,
	}
}

// IntOr is used to extract an environment variable into a Go int. If the
// environment variable is not set or is empty, then the alt value is used
// instead.
func IntOr(i *int, alt int) Parser {
	*i = alt
	return &intParser{
		required:    false,
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

// Float is used to extract an environment variable into a Go float64. If
// required is true, then an error is returned if the environment variable is
// not set or is empty.
func Float(f *float64, required bool) Parser {
	return &floatParser{
		required:    required,
		destination: f,
	}
}

// FloatOr is used to extract an environment variable into a Go float64. If
// the environment variable is not set or is emty, then the alt value is used
// instead.
func FloatOr(f *float64, alt float64) Parser {
	*f = alt
	return &floatParser{
		required:    false,
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

// Bool is used to extract an environment variable into a Go bool. If required
// is true, then an error is returned if the environment variable is not set or
// is empty.
func Bool(b *bool, required bool) Parser {
	return &boolParser{
		required:    required,
		destination: b,
	}
}

// BoolOr is used to extract an environment variable into a Go bool. If the
// environment variable is not set or is empty, then the alt value is used
// instead.
func BoolOr(b *bool, alt bool) Parser {
	*b = alt
	return &boolParser{
		required:    false,
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
var OS Environment = new(osEnv)

// File is an implementation of Environment that reads environment variables
// from a file.
//
// e.g. /etc/os-release
func File(filename string) Environment {
	return &fileEnv{filename: filename}
}

type fileEnv struct {
	filename string
}

func (e *fileEnv) Getenv(key string) string {
	f, err := os.Open(e.filename)
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		idx := strings.Index(line, "=")
		if idx < 1 || idx >= len(line)-1 {
			continue
		}
		if line[0:idx] == key {
			return line[idx+1:]
		}
	}

	return ""
}

// Map is an implementation of Environment that uses a given map[string]string
// to emulate a set of environment variables. Useful for testing.
//
//	m := Map(map[string]string {...})
//
//	func f(e env.Environment) {
//	  e.Getenv("FOOBAR")
//	}
func Map(m map[string]string) Environment {
	return &mapEnv{m: maps.Clone(m)}
}

type mapEnv struct {
	m map[string]string
}

func (m *mapEnv) Getenv(key string) string {
	return m.m[key]
}
