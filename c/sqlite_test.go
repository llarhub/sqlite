package sqlite

import (
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/llarhub/sqlite"
)

func TestInMemoryQuery(t *testing.T) {
	var db *sqlite.Sqlite3
	if rc := sqlite.Open(c.Str(":memory:"), &db); rc != sqlite.OK || db == nil {
		t.Fatalf("Open rc=%d db=%p", rc, db)
	}
	defer db.Close()

	var stmt *sqlite.Stmt
	if rc := db.PrepareV2(c.Str("SELECT ?1 + 1, upper(?2)"), -1, &stmt, nil); rc != sqlite.OK || stmt == nil {
		t.Fatalf("PrepareV2 rc=%d: %s", rc, c.GoString(db.Errmsg()))
	}
	defer stmt.Finalize()
	if rc := stmt.BindInt(1, 41); rc != sqlite.OK {
		t.Fatalf("BindInt rc=%d", rc)
	}
	if rc := stmt.BindText(2, c.Str("sqlite"), -1, nil); rc != sqlite.OK {
		t.Fatalf("BindText rc=%d", rc)
	}
	if rc := stmt.Step(); rc != sqlite.ROW {
		t.Fatalf("Step rc=%d: %s", rc, c.GoString(db.Errmsg()))
	}
	if got := stmt.ColumnInt(0); got != 42 {
		t.Fatalf("ColumnInt=%d, want 42", got)
	}
	got := stmt.ColumnText(1)
	if got == nil || c.GoString((*c.Char)(unsafe.Pointer(got))) != "SQLITE" {
		t.Fatalf("ColumnText=%q, want SQLITE", c.GoString((*c.Char)(unsafe.Pointer(got))))
	}
	if rc := stmt.Step(); rc != sqlite.DONE {
		t.Fatalf("second Step rc=%d", rc)
	}
}

func TestExecAndTransaction(t *testing.T) {
	var db *sqlite.Sqlite3
	if rc := sqlite.Open(c.Str(":memory:"), &db); rc != sqlite.OK || db == nil {
		t.Fatalf("Open rc=%d", rc)
	}
	defer db.Close()

	for _, sql := range []string{
		"CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT)",
		"BEGIN",
		"INSERT INTO items(name) VALUES ('one'), ('two')",
		"COMMIT",
	} {
		if rc := db.Exec(c.Str(sql), nil, nil, nil); rc != sqlite.OK {
			t.Fatalf("Exec(%q) rc=%d: %s", sql, rc, c.GoString(db.Errmsg()))
		}
	}
	if got := db.Changes(); got != 2 {
		t.Fatalf("Changes=%d, want 2", got)
	}
}

func TestApplicationFunction(t *testing.T) {
	var db *sqlite.Sqlite3
	if rc := sqlite.Open(c.Str(":memory:"), &db); rc != sqlite.OK || db == nil {
		t.Fatalf("Open rc=%d", rc)
	}
	defer db.Close()

	addOne := func(ctx *sqlite.Context, argc c.Int, argv **sqlite.Value) {
		if argc != 1 {
			ctx.ResultError(c.Str("expected one argument"), -1)
			return
		}
		ctx.ResultInt(c.Index(argv, 0).Int() + 1)
	}
	if rc := db.CreateFunction(c.Str("addone"), 1, sqlite.UTF8|sqlite.DETERMINISTIC, nil, addOne, nil, nil); rc != sqlite.OK {
		t.Fatalf("CreateFunction rc=%d: %s", rc, c.GoString(db.Errmsg()))
	}

	var stmt *sqlite.Stmt
	if rc := db.PrepareV2(c.Str("SELECT addone(41)"), -1, &stmt, nil); rc != sqlite.OK {
		t.Fatalf("PrepareV2 rc=%d: %s", rc, c.GoString(db.Errmsg()))
	}
	defer stmt.Finalize()
	if rc := stmt.Step(); rc != sqlite.ROW || stmt.ColumnInt(0) != 42 {
		t.Fatalf("function result rc=%d value=%d", rc, stmt.ColumnInt(0))
	}
}
