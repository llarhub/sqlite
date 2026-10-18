# SQLite

[![Go Reference](https://pkg.go.dev/badge/github.com/llarhub/sqlite.svg)](https://pkg.go.dev/github.com/llarhub/sqlite)
[![GitHub release](https://img.shields.io/github/v/tag/llarhub/sqlite.svg?label=release)](https://github.com/llarhub/sqlite/releases)
[![LLGo](https://img.shields.io/badge/powered_by-LLGo-green.svg)](https://github.com/xgo-dev/llgo)
[![XGo](https://img.shields.io/badge/project-XGo-blue.svg)](https://github.com/goplus/xgo)

LLGo bindings for SQLite, generated from the SQLite3 C API. The package exposes database connections, prepared statements, transactions, values, result callbacks, and application-defined SQL functions for XGo and Go programs compiled with LLGo.

## Installation

Install [LLGo](https://github.com/xgo-dev/llgo), SQLite, and `pkg-config`, then add the module with:

```bash
go get github.com/llarhub/sqlite
```

### macOS (Homebrew)

```bash
brew install sqlite pkg-config
```

### Linux (Debian or Ubuntu)

```bash
sudo apt update
sudo apt install libsqlite3-dev pkg-config
```

The binding links with `pkg-config --libs sqlite3` and `-lsqlite3`. Use the SQLite development package for your platform so that the linker can find `libsqlite3`.

## Usage

Open an in-memory database, bind a parameter, and read a result from a prepared statement:

```go
package main

import (
	"github.com/goplus/lib/c"
	"github.com/llarhub/sqlite"
)

func main() {
	var db *sqlite.Sqlite3
	if sqlite.Open(c.Str(":memory:"), &db) != sqlite.OK {
		panic("could not open database")
	}
	defer db.Close()

	var stmt *sqlite.Stmt
	if db.PrepareV2(c.Str("SELECT ?1 + 1"), -1, &stmt, nil) != sqlite.OK {
		panic(c.GoString(db.Errmsg()))
	}
	defer stmt.Finalize()

	stmt.BindInt(1, 41)
	if stmt.Step() == sqlite.ROW {
		println(stmt.ColumnInt(0)) // 42
	}
}
```

Use `Sqlite3.Exec` for schema and transaction statements, and `Stmt.Bind*`, `Stmt.Step`, `Stmt.Column*`, and `Stmt.Finalize` for prepared queries. `Sqlite3.CreateFunction` and the `Context.Result*` methods expose SQLite application-defined SQL functions.
