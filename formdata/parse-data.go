package formdata // import "go.gophers.dev/pkgs/extractors/formdata"

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrNonSingleValue  = errors.New("expected single value")
	ErrFieldNotPresent = errors.New("requested field does not exist")
	ErrParseFailure    = errors.New("could not parse value")
)

func Parse(data url.Values, schema Schema) error {
	for name, parser := range schema {
		values, exists := data[name]
		if !exists {
			return ErrFieldNotPresent
		}

		if err := parser.Parse(values); err != nil {
			return errors.Wrap(err, ErrParseFailure.Error())
		}
	}
	return nil
}

func ParseForm(r *http.Request, schema Schema) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return Parse(r.PostForm, schema)
}

// A Schema describes how a set of url.Values should be parsed.
// Typically these are coming from an http.Request.Form from inside an
// http.Handler responding to an inbound request.
type Schema map[string]Parser

// do we care about multi-value? we could provide parsers into slices
// automatically, for example

type Parser interface {
	Parse([]string) error
}

type stringParser struct {
	destination *string
}

func String(s *string) Parser {
	return &stringParser{destination: s}
}

func (p *stringParser) Parse(values []string) error {
	if len(values) != 1 {
		return ErrNonSingleValue
	}

	*p.destination = values[0]
	return nil
}

type intParser struct {
	destination *int
}

func Int(i *int) Parser {
	return &intParser{destination: i}
}

func (p *intParser) Parse(values []string) error {
	if len(values) != 1 {
		return ErrNonSingleValue
	}

	i, err := strconv.Atoi(values[0])
	if err != nil {
		return err
	}

	*p.destination = i
	return nil
}

type floatParser struct {
	destination *float64
}

func Float(f *float64) Parser {
	return &floatParser{destination: f}
}

func (p *floatParser) Parse(values []string) error {
	if len(values) != 1 {
		return ErrNonSingleValue
	}

	f, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return err
	}

	*p.destination = f
	return nil
}

type boolParser struct {
	destination *bool
}

func Bool(b *bool) Parser {
	return &boolParser{destination: b}
}

func (p *boolParser) Parse(values []string) error {
	if len(values) != 1 {
		return ErrNonSingleValue
	}

	b, err := strconv.ParseBool(values[0])
	if err != nil {
		return err
	}

	*p.destination = b
	return nil
}

// todo: slice variants?
