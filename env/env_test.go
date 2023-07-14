package env

import (
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
