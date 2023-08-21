extractors
==========

The `extractors` module provides libraries for defining a schema to easily and safely extract values from environment variables,
URL path elements, and HTML form values.

![GitHub](https://img.shields.io/github/license/shoenig/extractors.svg)
[![Run CI Tests](https://github.com/shoenig/extractors/actions/workflows/ci.yml/badge.svg)](https://github.com/shoenig/extractors/actions/workflows/ci.yml)


# Getting Started

The `extractors` package can be installed by running

```shell-session
go get github.com/shoenig/extractors@latest
```

```go
import (
    github.com/shoenig/extractors/env      // extract values from environment variables
    github.com/shoenig/extractors/urlpath  // extract elements from url paths
    github.com/shoenig/extractors/formdata // extract values from html data
)
```

#### env example

Use the `env` package to parse values from environment variables.

```go
import github.com/shoenig/go-conceal // for storing sensitive values
```

```go
var (
    go111module string
    sshPID      int
    password    *conceal.Text
)

_ = env.ParseOS(env.Schema{
    "GO111MODULE":   env.String(&go111module, false),
    "SSH_AGENT_PID": env.Int(&sshPID, true),
    "PASSWORD":      env.Secret(&password, true),
})
```

#### formdata example

Use the `formdata` package to parse values from `url.Values` (typically coming
from ``*http.Request.Form` objects from inbound requests.

```go
// typically coming from a *http.Request.Form
values := url.Values{
    "user":     []string{"bob"},
    "age":      []string{"45"},
}

var (
    user string
    age  int
)

_ = formdata.Parse(values, formdata.Schema{
    "user": formdata.String(&user),
    "age":  formdata.Int(&age),
})
```

#### urlpath example

Use the `urlpath` package to parse URL path elements when using a `gorilla/mux`
router.

```go
// with a mux handler definition like
router.Handle("/{kind}/{id}")

// in the handler implementation, parse the *http.Request URL with
var (
    kind string
    id   int
)

_ = urlpath.Parse(request, urlpath.Schema{
    "kind": urlpath.String(&kind),
    "id":   urlpath.Int(&id),
})
```

# Contributing

The `github.com/shoenig/extractors` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `github.com/shoenig/extractors` module is open source under the [BSD-3-Clause](LICENSE) license.
