package formdata

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/shoenig/test/must"
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
	must.NoError(t, err)
	must.EqCmp(t, "1", one)
	must.EqCmp(t, 2, two)
	must.EqCmp(t, 3.1, three)
	must.True(t, four)
}

func Test_Parse_HTMLForm(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "/", nil)
	must.NoError(t, err)

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
	must.NoError(t, err2)
	must.EqCmp(t, "1", one)
	must.EqCmp(t, 2, two)
	must.EqCmp(t, 3.1, three)
	must.True(t, four)
}

func Test_Parse_HTMLForm_not_ready(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "/", nil)
	must.NoError(t, err)

	var one string

	// not yet a valid form, never had the FormValues field set
	err2 := ParseForm(request, Schema{
		"one": String(&one),
	})
	must.Error(t, err2)
}

func Test_Parse_key_missing(t *testing.T) {
	data := url.Values{
		"one": []string{"1"},
	}

	var two int
	err := Parse(data, Schema{
		"two": Int(&two),
	})
	must.Error(t, err)
}

func Test_Parse_string_value_missing(t *testing.T) {
	data := url.Values{
		"one": []string{},
	}

	var one string
	err := Parse(data, Schema{
		"one": String(&one),
	})
	must.Error(t, err)
}

func Test_Parse_int_value_missing(t *testing.T) {
	data := url.Values{
		"two": []string{},
	}

	var two int
	err := Parse(data, Schema{
		"two": Int(&two),
	})
	must.Error(t, err)
}

func Test_Parse_int_malformed(t *testing.T) {
	data := url.Values{
		"two": []string{"not an int"},
	}

	var two int
	err := Parse(data, Schema{
		"two": Int(&two),
	})
	must.Error(t, err)
}

func Test_Parse_float_value_missing(t *testing.T) {
	data := url.Values{
		"three": []string{},
	}

	var three float64
	err := Parse(data, Schema{
		"three": Float(&three),
	})
	must.Error(t, err)
}

func Test_Parse_float_malformed(t *testing.T) {
	data := url.Values{
		"three": []string{"not a float"},
	}

	var three float64
	err := Parse(data, Schema{
		"three": Float(&three),
	})
	must.Error(t, err)
}

func Test_Parse_bool_value_missing(t *testing.T) {
	data := url.Values{
		"four": []string{},
	}

	var four bool
	err := Parse(data, Schema{
		"four": Bool(&four),
	})
	must.Error(t, err)
}

func Test_Parse_bool_malformed(t *testing.T) {
	data := url.Values{
		"four": []string{"not a bool"},
	}

	var four bool
	err := Parse(data, Schema{
		"four": Bool(&four),
	})
	must.Error(t, err)
}
