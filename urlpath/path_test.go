// Copyright (c) Seth Hoenig
// SPDX-License-Identifier: BSD-3-Clause

package urlpath

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/shoenig/test/must"
)

func Test_Parse(t *testing.T) {
	router := mux.NewRouter()
	executed := false

	router.HandleFunc("/v1/{foo}/{bar}", func(_ http.ResponseWriter, r *http.Request) {
		var foo string
		var bar int

		err := Parse(r, Schema{
			"foo": String(&foo),
			"bar": Int(&bar),
		})

		must.NoError(t, err)
		must.EqOp(t, "blah", foo)
		must.EqOp(t, 31, bar)
		executed = true
	})

	w := httptest.NewRecorder()
	ctx := context.Background()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v1/blah/31", nil)
	must.NoError(t, err)

	router.ServeHTTP(w, request)
	must.True(t, executed)
}

func Test_ParseValues(t *testing.T) {
	var foo string
	var bar int
	var id uint64

	values := map[string]string{
		"foo": "blah",
		"bar": "21",
		"id":  "42",
	}

	err := ParseValues(values, Schema{
		"foo": String(&foo),
		"bar": Int(&bar),
		"id":  UInt64(&id),
	})

	must.NoError(t, err)
	must.EqOp(t, "blah", foo)
	must.EqOp(t, 21, bar)
	must.EqOp(t, 42, id)
}

func Test_ParseValues_incompatible(t *testing.T) {
	var foo string
	var bar int

	values := map[string]string{
		"foo": "blah",
		"bar": "not an int",
	}

	err := ParseValues(values, Schema{
		"foo": String(&foo),
		"bar": Int(&bar),
	})

	must.Error(t, err)
}

func Test_ParseValues_missing(t *testing.T) {
	var foo string
	var bar int

	values := map[string]string{
		"foo": "blah",
	}

	err := ParseValues(values, Schema{
		"foo": String(&foo),
		"bar": Int(&bar),
	})

	must.Error(t, err)
}

func Test_Parameter_String(t *testing.T) {
	p := Parameter("foo")
	s := p.String()
	must.EqOp(t, "{foo}", s)
}
