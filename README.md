authz [![GoDoc](https://godoc.org/github.com/casbin/chi-authz?status.svg)](https://godoc.org/github.com/casbin/chi-authz)
======

chi-authz is an authorization middleware for [Chi](https://github.com/pressly/chi), it's based on [https://github.com/casbin/casbin](https://github.com/casbin/casbin).

## Installation

    go get github.com/casbin/chi-authz

## Simple Example

```Go
package main

import (
	"net/http"

	"github.com/casbin/chi-authz"
	"github.com/casbin/casbin"
	"github.com/pressly/chi"
)

func main() {
	router := chi.NewRouter()

	// load the casbin model and policy from files, database is also supported.
	e := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")
	router.Use(authz.Authorizer(e))

	// define your handler, this is just an example to return HTTP 200 for any requests.
	// the access that is denied by authz will return HTTP 403 error.
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}
```

## Getting Help

- [casbin](https://github.com/casbin/casbin)

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
