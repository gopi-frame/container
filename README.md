# container
[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/container.svg)](https://pkg.go.dev/github.com/gopi-frame/container)
[![Go](https://github.com/gopi-frame/container/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/gopi-frame/container/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gopi-frame/container/graph/badge.svg?token=UGVGP6QF5O)](https://codecov.io/gh/gopi-frame/container)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/container)](https://goreportcard.com/report/github.com/gopi-frame/container)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package container provides an implementation of [Container](https://github.com/gopi-frame/contract/container).

## Installation

```shell
go get -u github.com/gopi-frame/container
```

## Import

```go
import "github.com/gopi-frame/container"
```

## Usage

```go
package main

import (
    "database/sql"
    "github.com/gopi-frame/container"
)

var db *sql.DB
var err error

func init() {
    db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
    if err != nil {
        panic(err)
    }
}

func main() {
    var c = container.New[*sql.DB]()
    c.Set("db", db)
    c.Lazy("db2", func() (*sql.DB, error) {
        return sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test2")
    })
    db := c.Get("db")
    // do something with db
    db2 := c.Get("db2")
    // do something with db2
    newDB2, err := c.Make("db2")
    if err != nil {
        panic(err)
    }
    // do something with newDB2
}
```
