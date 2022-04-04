package main

import (
	"TinyORM/session"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "tiny.db")
	if err != nil {
		t.Fatal("failed to connect")
	}
	return engine
}

type User struct {
	Name string `tinyorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(session *session.Session) (result interface{}, err error) {
		_ = session.Model(&User{}).CreateTable()
		_, err = session.Insert(&User{"Tom", 18})
		// 故意抛出一个错误
		return nil, errors.New("error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, err := s.Raw("DROP TABLE IF EXISTS User;").Exec()
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if err != nil {
		t.Fatal(err)
	}
	_ = engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}