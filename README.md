extractors
==========

Extract values from text using typed schema

[![Go Report Card](https://goreportcard.com/badge/gophers.dev/pkgs/extractors)](https://goreportcard.com/report/gophers.dev/pkgs/extractors)
[![Build Status](https://travis-ci.com/shoenig/extractors.svg?branch=master)](https://travis-ci.com/shoenig/extractors)
[![GoDoc](https://godoc.org/gophers.dev/pkgs/extractors?status.svg)](https://godoc.org/gophers.dev/pkgs/extractors)
![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/shoenig/extractors.svg)
![GitHub](https://img.shields.io/github/license/shoenig/extractors.svg)

# Project Overview

Module `gophers.dev/pkgs/extractors` provides packages for extracting values
from sources of text. By providing a typed schema with named keys, values can
be parsed from text in a type-safe way.

# Getting Started

The `extractors` package can be installed by running
```
$ go get gophers.dev/pkgs/extractors
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

The `gophers.dev/pkgs/extractors` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `gophers.dev/pkgs/extractors` module is open source under the [BSD-3-Clause](LICENSE) license.
