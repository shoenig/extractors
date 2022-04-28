package env

import (
	"testing"

	"github.com/shoenig/test/must"
)

func Test_Parse_required(t *testing.T) {
	env := NewEnvironmentMock(t)
	defer env.MinimockFinish()

	env.GetenvMock.When("FOO").Then("foo")
	env.GetenvMock.When("BAR").Then("12")
	env.GetenvMock.When("BAZ").Then("3.14")
	env.GetenvMock.When("B1").Then("true")
	env.GetenvMock.When("B2").Then("1")

	var (
		foo string
		bar int
		baz float64
		b1  bool
		b2  bool
	)

	err := Parse(env, Schema{
		"FOO": String(&foo, true),
		"BAR": Int(&bar, true),
		"BAZ": Float(&baz, true),
		"B1":  Bool(&b1, true),
		"B2":  Bool(&b2, true),
	})

	must.NoError(t, err)
	must.EqCmp(t, "foo", foo)
	must.EqCmp(t, 12, bar)
	must.EqCmp(t, 3.14, baz)
	must.True(t, b1)
	must.True(t, b2)
}

func Test_Parse_optional(t *testing.T) {
	env := NewEnvironmentMock(t)
	defer env.MinimockFinish()

	env.GetenvMock.When("FOO").Then("")
	env.GetenvMock.When("BAR").Then("")
	env.GetenvMock.When("BAZ").Then("")
	env.GetenvMock.When("B1").Then("")

	var (
		foo string
		bar int
		baz float64
		b1  bool
	)

	err := Parse(env, Schema{
		"FOO": String(&foo, false),
		"BAR": Int(&bar, false),
		"BAZ": Float(&baz, false),
		"B1":  Bool(&b1, false),
	})

	must.NoError(t, err)
	must.EqCmp(t, "", foo)
	must.EqCmp(t, 0, bar)
	must.EqCmp(t, 0.0, baz)
	must.False(t, b1)
}

func Test_Parse_required_missing(t *testing.T) {
	env := NewEnvironmentMock(t)
	defer env.MinimockFinish()

	env.GetenvMock.When("FOO").Then("")
	env.GetenvMock.When("BAR").Then("")
	env.GetenvMock.When("BAZ").Then("")
	env.GetenvMock.When("B1").Then("")

	{
		var foo string
		err := Parse(env, Schema{
			"FOO": String(&foo, true),
		})
		must.Error(t, err)
	}

	{
		var bar int
		err := Parse(env, Schema{
			"BAR": Int(&bar, true),
		})
		must.Error(t, err)
	}

	{
		var baz float64
		err := Parse(env, Schema{
			"BAZ": Float(&baz, true),
		})
		must.Error(t, err)
	}

	{
		var b1 bool
		err := Parse(env, Schema{
			"B1": Bool(&b1, true),
		})
		must.Error(t, err)
	}
}

func Test_Parse_fail(t *testing.T) {
	env := NewEnvironmentMock(t)
	defer env.MinimockFinish()

	env.GetenvMock.When("BAR").Then("abc")
	env.GetenvMock.When("BAZ").Then("abc")

	{
		var bar int
		err := Parse(env, Schema{
			"BAR": Int(&bar, true),
		})
		must.Error(t, err)
	}

	{
		var baz float64
		err := Parse(env, Schema{
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
