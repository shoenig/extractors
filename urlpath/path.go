// Copyright (c) Seth Hoenig
// SPDX-License-Identifier: BSD-3-Clause

package urlpath

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Typical usage:
//
//    // in the mux definition, e.g.
//    router.Handle("/v1/{foo}/{bar}, newHandler())
//
//    // in the handler implementation, e.g.
//    var foo string
//    var bar int
//    err := urlpath.Parse(request, urlpath.Schema{
//        "foo": urlpath.String(&foo),
//        "bar": urlpath.Int(&bar),
//    })

// A Parameter is a named element of a URL route,
// encoded such that a gorilla router interprets it
// as a path parameter.
type Parameter string

func (p Parameter) String() string {
	return "{" + string(p) + "}"
}

// Name returns the name of the parameter.
func (p Parameter) Name() string {
	return string(p)
}

// A Schema describes how path variables should be parsed.
// Currently only int and string types are supported.
type Schema map[Parameter]Parser

// Parse will parse the URL path vars from r given the
// element names and parsers defined in schema.
//
// This method only works with requests being processed by
// handlers of a gorilla/mux.
func Parse(r *http.Request, schema Schema) error {
	return ParseValues(mux.Vars(r), schema)
}

// ParseValues will parse the parameters in vars given the
// element names and parsers defined in schema.
//
// Most use cases will be parsing values coming from an *http.Request,
// which can be done conveniently with Parse.
func ParseValues(values map[string]string, schema Schema) error {
	for name, parser := range schema {
		value, exists := values[name.Name()]
		if !exists {
			return fmt.Errorf("url path element not present: %q", name)
		}

		if err := parser.Parse(value); err != nil {
			return fmt.Errorf("could not parse url path variable: %w", err)
		}
	}
	return nil
}

// A Parser parses raw input into a destination variable.
type Parser interface {
	Parse(string) error
}

type stringParser struct {
	destination *string
}

// String creates a parser that will parse a path element into s.
func String(s *string) Parser {
	return &stringParser{destination: s}
}

func (p *stringParser) Parse(s string) error {
	*p.destination = s
	return nil
}

type intParser struct {
	destination *int
}

// Int creates a Parser that will parse a path element into i.
func Int(i *int) Parser {
	return &intParser{destination: i}
}

func (p *intParser) Parse(s string) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*p.destination = i
	return nil
}

type uint64Parser struct {
	destination *uint64
}

func UInt64(i *uint64) Parser {
	return &uint64Parser{destination: i}
}

func (u uint64Parser) Parse(s string) error {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*u.destination = i
	return nil
}
