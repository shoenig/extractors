package env

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/go-conceal"
	"github.com/shoenig/test/must"
)

func Test_Parse_required(t *testing.T) {
	t.Setenv("FOO", "foo")
	t.Setenv("BAR", "12")
	t.Setenv("BAZ", "3.14")
	t.Setenv("B1", "true")
	t.Setenv("B2", "1")
	t.Setenv("PASSWORD", "hunter2")

	var (
		foo  string
		bar  int
		baz  float64
		b1   bool
		b2   bool
		pass *conceal.Text
	)

	err := Parse(OS, Schema{
		"FOO":      String(&foo, true),
		"BAR":      Int(&bar, true),
		"BAZ":      Float(&baz, true),
		"B1":       Bool(&b1, true),
		"B2":       Bool(&b2, true),
		"PASSWORD": Secret(&pass, true),
	})

	must.NoError(t, err)
	must.Eq(t, "foo", foo)
	must.Eq(t, 12, bar)
	must.Eq(t, 3.14, baz)
	must.True(t, b1)
	must.True(t, b2)
	must.Eq(t, "hunter2", pass.Unveil())
}

func Test_ParseOr(t *testing.T) {
	t.Setenv("FOO", "foo")
	t.Setenv("BAR", "")
	t.Setenv("I1", "42")
	t.Setenv("I2", "")
	t.Setenv("F1", "1.1")
	t.Setenv("F2", "")
	t.Setenv("B1", "true")
	t.Setenv("B2", "")

	var (
		foo string
		bar string
		i1  int
		i2  int
		f1  float64
		f2  float64
		b1  bool
		b2  bool
	)

	err := Parse(OS, Schema{
		"FOO": StringOr(&foo, "baz"),
		"BAR": StringOr(&bar, "baz"),
		"I1":  IntOr(&i1, 77),
		"I2":  IntOr(&i2, 77),
		"F1":  FloatOr(&f1, 50.5),
		"F2":  FloatOr(&f2, 50.5),
		"B1":  BoolOr(&b1, false),
		"B2":  BoolOr(&b2, true),
	})

	must.NoError(t, err)
	must.Eq(t, "foo", foo)
	must.Eq(t, "baz", bar)
	must.Eq(t, 42, i1)
	must.Eq(t, 77, i2)
	must.Eq(t, f1, 1.1)
	must.Eq(t, f2, 50.5)
	must.True(t, b1)
	must.True(t, b2)
}

func Test_Parse_optional(t *testing.T) {
	t.Setenv("FOO", "")
	t.Setenv("BAR", "")
	t.Setenv("BAZ", "")
	t.Setenv("B1", "")
	t.Setenv("PASSWORD", "")

	var (
		foo  string
		bar  int
		baz  float64
		b1   bool
		pass *conceal.Text
	)

	err := Parse(OS, Schema{
		"FOO":      String(&foo, false),
		"BAR":      Int(&bar, false),
		"BAZ":      Float(&baz, false),
		"B1":       Bool(&b1, false),
		"PASSWORD": Secret(&pass, false),
	})

	must.NoError(t, err)
	must.Eq(t, "", foo)
	must.Eq(t, 0, bar)
	must.Eq(t, 0.0, baz)
	must.False(t, b1)
	must.Nil(t, pass)
}

func Test_Parse_required_missing(t *testing.T) {
	t.Setenv("FOO", "")
	t.Setenv("BAR", "")
	t.Setenv("BAZ", "")
	t.Setenv("B1", "")
	t.Setenv("PASSWORD", "")
	{
		var foo string
		err := Parse(OS, Schema{
			"FOO": String(&foo, true),
		})
		must.Error(t, err)
	}

	{
		var bar int
		err := Parse(OS, Schema{
			"BAR": Int(&bar, true),
		})
		must.Error(t, err)
	}

	{
		var baz float64
		err := Parse(OS, Schema{
			"BAZ": Float(&baz, true),
		})
		must.Error(t, err)
	}

	{
		var b1 bool
		err := Parse(OS, Schema{
			"B1": Bool(&b1, true),
		})
		must.Error(t, err)
	}
	{
		var cred *conceal.Text
		err := Parse(OS, Schema{
			"PASSWORD": Secret(&cred, true),
		})
		must.Error(t, err)
		must.Nil(t, cred)
	}
}

func Test_Parse_fail(t *testing.T) {
	t.Setenv("BAR", "abc")
	t.Setenv("BAZ", "abc")

	{
		var bar int
		err := Parse(OS, Schema{
			"BAR": Int(&bar, true),
		})
		must.Error(t, err)
	}

	{
		var baz float64
		err := Parse(OS, Schema{
			"BAZ": Float(&baz, true),
		})
		must.Error(t, err)
	}
}

func Test_ParseOS(t *testing.T) {
	var xTerm string

	err := ParseOS(Schema{
		"XTERM": String(&xTerm, false),
	})

	must.NoError(t, err)
	t.Log("xterm value:", xTerm)
}

func Test_File(t *testing.T) {
	temp := filepath.Join(t.TempDir(), "test.env")
	file, err := os.OpenFile(temp, os.O_CREATE|os.O_WRONLY, 0644)
	must.NoError(t, err)

	text := `
ONE=1
TWO=two
THREE=three
`

	_, err = io.WriteString(file, text)
	must.NoError(t, err)

	f := File(temp)

	one := f.Getenv("ONE")
	must.Eq(t, "1", one)

	two := f.Getenv("TWO")
	must.Eq(t, "two", two)

	three := f.Getenv("THREE")
	must.Eq(t, "three", three)

	missing := f.Getenv("FOUR")
	must.Eq(t, "", missing)
}

func Test_ParseFile(t *testing.T) {
	temp := filepath.Join(t.TempDir(), "test.env")
	file, err := os.OpenFile(temp, os.O_CREATE|os.O_WRONLY, 0644)
	must.NoError(t, err)

	text := `
ONE=1
TWO=two
THREE=three
`

	_, err = io.WriteString(file, text)
	must.NoError(t, err)

	var (
		one   int
		two   string
		three string
	)

	err = ParseFile(temp, Schema{
		"ONE":   Int(&one, true),
		"TWO":   String(&two, true),
		"THREE": String(&three, true),
	})

	must.NoError(t, err)
	must.Eq(t, 1, one)
	must.Eq(t, "two", two)
	must.Eq(t, "three", three)
}
