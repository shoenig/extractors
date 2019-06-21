package formdata

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse_singles(t *testing.T) {
	data := url.Values{
		"one":   []string{"1"},
		"two":   []string{"2"},
		"three": []string{"3.1"},
		"four":  []string{"true"},
	}

	var (
		one   string
		two   int
		three float64
		four  bool
	)

	err := Parse(data, Schema{
		"one":   String(&one),
		"two":   Int(&two),
		"three": Float(&three),
		"four":  Bool(&four),
	})
	require.NoError(t, err)
	require.Equal(t, "1", one)
	require.Equal(t, 2, two)
	require.Equal(t, 3.1, three)
	require.Equal(t, true, four)
}

func Test_Parse_HTMLForm(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "/", nil)
	require.NoError(t, err)

	request.PostForm = make(url.Values)
	request.PostForm.Set("one", "1")
	request.PostForm.Set("two", "2")
	request.PostForm.Set("three", "3.1")
	request.PostForm.Set("four", "true")

	var (
		one   string
		two   int
		three float64
		four  bool
	)

	err2 := ParseForm(request, Schema{
		"one":   String(&one),
		"two":   Int(&two),
		"three": Float(&three),
		"four":  Bool(&four),
	})
	require.NoError(t, err2)
	require.Equal(t, "1", one)
	require.Equal(t, 2, two)
	require.Equal(t, 3.1, three)
	require.Equal(t, true, four)
}

func Test_Parse_HTMLForm_not_ready(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "/", nil)
	require.NoError(t, err)

	var one string

	// not yet a valid form, never had the FormValues field set
	err2 := ParseForm(request, Schema{
		"one": String(&one),
	})
	require.Error(t, err2)

}

func Test_Parse_key_missing(t *testing.T) {
	data := url.Values{
		"one": []string{"1"},
	}

	var two int
	err := Parse(data, Schema{
		"two": Int(&two),
	})
	require.Error(t, err)
}

func Test_Parse_string_value_missing(t *testing.T) {
	data := url.Values{
		"one": []string{},
	}

	var one string
	err := Parse(data, Schema{
		"one": String(&one),
	})
	require.Error(t, err)
}

func Test_Parse_int_value_missing(t *testing.T) {
	data := url.Values{
		"two": []string{},
	}

	var two int
	err := Parse(data, Schema{
		"two": Int(&two),
	})
	require.Error(t, err)
}

func Test_Parse_int_malformed(t *testing.T) {
	data := url.Values{
		"two": []string{"not an int"},
	}

	var two int
	err := Parse(data, Schema{
		"two": Int(&two),
	})
	require.Error(t, err)
}

func Test_Parse_float_value_missing(t *testing.T) {
	data := url.Values{
		"three": []string{},
	}

	var three float64
	err := Parse(data, Schema{
		"three": Float(&three),
	})
	require.Error(t, err)
}

func Test_Parse_float_malformed(t *testing.T) {
	data := url.Values{
		"three": []string{"not a float"},
	}

	var three float64
	err := Parse(data, Schema{
		"three": Float(&three),
	})
	require.Error(t, err)
}

func Test_Parse_bool_value_missing(t *testing.T) {
	data := url.Values{
		"four": []string{},
	}

	var four bool
	err := Parse(data, Schema{
		"four": Bool(&four),
	})
	require.Error(t, err)
}

func Test_Parse_bool_malformed(t *testing.T) {
	data := url.Values{
		"four": []string{"not a bool"},
	}

	var four bool
	err := Parse(data, Schema{
		"four": Bool(&four),
	})
	require.Error(t, err)
}
