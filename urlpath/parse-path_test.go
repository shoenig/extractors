package urlpath

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	router := mux.NewRouter()
	executed := false

	router.HandleFunc("/v1/{foo}/{bar}", func(w http.ResponseWriter, r *http.Request) {
		var foo string
		var bar int

		err := Parse(r, Schema{
			"foo": String(&foo),
			"bar": Int(&bar),
		})

		require.NoError(t, err)
		require.Equal(t, "blah", foo)
		require.Equal(t, 31, bar)
		executed = true
	})

	w := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/v1/blah/31", nil)
	require.NoError(t, err)

	router.ServeHTTP(w, request)
	require.True(t, executed)
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

	require.NoError(t, err)
	require.Equal(t, "blah", foo)
	require.Equal(t, 21, bar)
	require.Equal(t, uint64(42), id)
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

	require.Error(t, err)
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

	require.Error(t, err)
}

func Test_Parameter_String(t *testing.T) {
	p := Parameter("foo")
	s := p.String()
	require.Equal(t, "{foo}", s)
}
