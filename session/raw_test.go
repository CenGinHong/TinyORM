package session

import (
	"TinyORM/dialect"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

type User struct {
	Name string `tinyorm:"PRIMARY KEY"`
	Age  int
}

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../tiny.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}
