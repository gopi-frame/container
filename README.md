# container

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
