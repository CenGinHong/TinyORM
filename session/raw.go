package session

import (
	"TinyORM/log"
	"database/sql"
	"strings"
)

type Session struct {
	db      *sql.DB         // 连接数据库的指针
	sql     strings.Builder // 拼接sql
	sqlVars []interface{}   // sql 需要填入的变量
}

func New(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec 执行sql
func (s *Session) Exec() (result sql.Result, err error) {
	// 执行后清空sql内容
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// 使用原生db执行
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
		return
	}
	return
}

// QueryRow 查询
func (s *Session) QueryRow() *sql.Row {
	// 执行后清空sql内容
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// 使用原生db执行，QueryRow只返回一行
	return s.DB().QueryRow(s.sql.String(), s.sqlVars)
}

// QueryRows 查询多行
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	// 执行后清空sql内容
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// Query返回多行
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}