extractors
==========

Extract values from text using typed schema

![GitHub](https://img.shields.io/github/license/shoenig/extractors.svg)

# Project Overview

Module `github.com/shoenig/extractors` provides packages for extracting values
from sources of text. By providing a typed schema with named keys, values can
be parsed from text in a type-safe way.

# Getting Started

The `extractors` package can be installed by running
```
$ go get github.com/shoenig/extractors
```

#### Example Usage of env
Use the `env` package to parse values from environment variables.
```golang
var (
    go111module string
    sshPID      int
)

_ = env.ParseOS(env.Schema{
    "GO111MODULE":   env.String(&go111module, false),
    "SSH_AGENT_PID": env.Int(&sshPID, true),
})

```

#### Example usage of formdata
Use the `formdata` package to parse values from `url.Values` (typically coming
from ``*http.Request.Form` objects from inbound requests.
```golang
// typically coming from a *http.Request.Form
values := url.Values{
    "user": []string{"bob"},
    "age":  []string{"45"},
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

#### Example usage of urlpath
Use the `urlpath` package to parse URL path elements when using a `gorilla/mux`
router.
```golang
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
