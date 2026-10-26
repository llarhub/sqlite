package sqlite

import (
	"github.com/goplus/lib/c"
)

// Opening A New Database Connection
// filename: Database filename (UTF-8)
func Open(filename *c.Char) (db *Sqlite3, err c.Int) {
	err = doOpen(filename, &db)
	return
}

// Opening A New Database Connection
// filename: Database filename (UTF-8)
// zVfs: Name of VFS module to use
func OpenV2(filename *c.Char, flags c.Int, zVfs *c.Char) (db *Sqlite3, err c.Int) {
	err = doOpenV2(filename, &db, flags, zVfs)
	return
}

// Compiling An SQL Statement
// tail: Pointer to unused portion of sql
func (self *Sqlite3) Prepare(sql string, tail **c.Char) (stmt *Stmt, err c.Int) {
	err = self.doPrepare(c.GoStringData(sql), c.Int(len(sql)), &stmt, tail)
	return
}

func (self *Sqlite3) PrepareV2(sql string, tail **c.Char) (stmt *Stmt, err c.Int) {
	err = self.doPrepareV2(c.GoStringData(sql), c.Int(len(sql)), &stmt, tail)
	return
}

func (self *Sqlite3) PrepareV3(sql string, flags c.Uint, tail **c.Char) (stmt *Stmt, err c.Int) {
	err = self.doPrepareV3(c.GoStringData(sql), c.Int(len(sql)), flags, &stmt, tail)
	return
}
