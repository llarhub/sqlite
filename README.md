# SQLite

[![Go Reference](https://pkg.go.dev/badge/github.com/llarhub/sqlite.svg)](https://pkg.go.dev/github.com/llarhub/sqlite)
[![GitHub release](https://img.shields.io/github/v/tag/llarhub/sqlite.svg?label=release)](https://github.com/llarhub/sqlite/releases)
[![LLGo](https://img.shields.io/badge/powered_by-LLGo-green.svg)](https://github.com/xgo-dev/llgo)
[![LLGo bindings](https://img.shields.io/badge/llarhub-LLGo_Bindings-brightgreen.svg)](https://github.com/xgo-dev/llgo/blob/main/doc/LLGo-Bindings.md)
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
	"github.com/goplus/lib/c/os"
	"github.com/llarhub/sqlite"
)

func main() {
	os.Remove(c.Str("test.db"))

	db, err := sqlite.Open(c.Str("test.db"))
	check(err, db, "sqlite: Open")
	defer db.Close()

	err = db.Exec(c.Str("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)"), nil, nil, nil)
	check(err, db, "sqlite: Exec CREATE TABLE")

	{
		stmt, err := db.PrepareV3("INSERT INTO users (id, name) VALUES (?, ?)", 0, nil)
		check(err, db, "sqlite: PrepareV3 INSERT")
		defer stmt.Finalize()

		stmt.BindInt(1, 100)
		stmt.BindText(2, c.Str("Hello World"), -1, nil)

		err = stmt.Step()
		checkDone(err, db, "sqlite: Step INSERT 1")

		stmt.Reset()
		stmt.BindInt(1, 200)
		stmt.BindText(2, c.Str("This is llgo"), -1, nil)

		err = stmt.Step()
		checkDone(err, db, "sqlite: Step INSERT 2")
	}
	{
		stmt, err := db.PrepareV3("SELECT * FROM users", 0, nil)
		check(err, db, "sqlite: PrepareV3 SELECT")
		defer stmt.Finalize()

		for stmt.Step() == sqlite.ROW {
			c.Printf(c.Str("==> id=%d, name=%s\n"), stmt.ColumnInt(0), stmt.ColumnText(1))
		}
		checkDone(db.Errcode(), db, "sqlite: Step done")
	}
}

func check(err c.Int, db *sqlite.Sqlite3, at string) {
	if err != sqlite.OK {
		c.Printf(c.Str("==> %s Error: (%d) %s\n"), c.AllocaCStr(at), err, db.Errmsg())
		c.Exit(1)
	}
}

func checkDone(err c.Int, db *sqlite.Sqlite3, at string) {
	if err != sqlite.DONE {
		check(err, db, at)
	}
}
```

Use `Sqlite3.Exec` for schema and transaction statements, and `Stmt.Bind*`, `Stmt.Step`, `Stmt.Column*`, and `Stmt.Finalize` for prepared queries. `Sqlite3.CreateFunction` and the `Context.Result*` methods expose SQLite application-defined SQL functions.
