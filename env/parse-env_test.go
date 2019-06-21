package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse_required(t *testing.T) {
	env := NewEnvironmentMock(t)
	defer env.MinimockFinish()

	env.GetenvMock.When("FOO").Then("foo")
	env.GetenvMock.When("BAR").Then("12")
	env.GetenvMock.When("BAZ").Then("3.14")

	var (
		foo string
		bar int
		baz float64
	)

	err := Parse(env, Schema{
		"FOO": String(&foo, true),
		"BAR": Int(&bar, true),
		"BAZ": Float(&baz, true),
	})

	require.NoError(t, err)
	require.Equal(t, "foo", foo)
	require.Equal(t, 12, bar)
	require.Equal(t, 3.14, baz)
}

func Test_Parse_optional(t *testing.T) {
	env := NewEnvironmentMock(t)
	defer env.MinimockFinish()

	env.GetenvMock.When("FOO").Then("")
	env.GetenvMock.When("BAR").Then("")
	env.GetenvMock.When("BAZ").Then("")

	var (
		foo string
		bar int
		baz float64
	)

	err := Parse(env, Schema{
		"FOO": String(&foo, false),
		"BAR": Int(&bar, false),
		"BAZ": Float(&baz, false),
	})

	require.NoError(t, err)
	require.Equal(t, "", foo)
	require.Equal(t, 0, bar)
	require.Equal(t, 0.0, baz)
}

func Test_Parse_required_missing(t *testing.T) {
	env := NewEnvironmentMock(t)
	defer env.MinimockFinish()

	env.GetenvMock.When("FOO").Then("")
	env.GetenvMock.When("BAR").Then("")
	env.GetenvMock.When("BAZ").Then("")

	{
		var foo string
		err := Parse(env, Schema{
			"FOO": String(&foo, true),
		})
		require.Error(t, err)
	}

	{
		var bar int
		err := Parse(env, Schema{
			"BAR": Int(&bar, true),
		})
		require.Error(t, err)
	}

	{
		var baz float64
		err := Parse(env, Schema{
			"BAZ": Float(&baz, true),
		})
		require.Error(t, err)
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
		require.Error(t, err)
	}

	{
		var baz float64
		err := Parse(env, Schema{
			"BAZ": Float(&baz, true),
		})
		require.Error(t, err)
	}
}

func Test_ParseOS(t *testing.T) {
	var xTerm string

	err := ParseOS(Schema{
		"XTERM": String(&xTerm, false),
	})

	require.NoError(t, err)
	t.Log("xterm value:", xTerm)
}
