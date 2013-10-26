package gorm

import (
	"database/sql"
	"errors"

	"strconv"
)

type Orm struct {
	TableName  string
	PrimaryKey string
	SqlResult  sql.Result
	Error      error
	Sql        string
	SqlVars    []interface{}
	model      *Model

	db          *sql.DB
	driver      string
	whereClause []map[string]interface{}
	selectStr   string
	orderStr    string
	offsetInt   int
	limitInt    int
	operation   string
}

func (s *Orm) Model(model interface{}) (err error) {
	s.model = s.toModel(model)
	s.TableName = s.model.TableName()
	s.PrimaryKey = s.model.PrimaryKeyDb()
	return
}

func (s *Orm) Where(querystring interface{}, args ...interface{}) *Orm {
	s.whereClause = append(s.whereClause, map[string]interface{}{"query": querystring, "args": args})
	return s
}

func (s *Orm) Limit(value interface{}) *Orm {
	switch value := value.(type) {
	case string:
		s.limitInt, _ = strconv.Atoi(value)
	case int:
		s.limitInt = value
	default:
		s.Error = errors.New("Can' understand the value of Limit, Should be int")
	}
	return s
}

func (s *Orm) Offset(value interface{}) *Orm {
	switch value := value.(type) {
	case string:
		s.offsetInt, _ = strconv.Atoi(value)
	case int:
		s.offsetInt = value
	default:
		s.Error = errors.New("Can' understand the value of Offset, Should be int")
	}
	return s
}

func (s *Orm) Order(value string, reorder ...bool) *Orm {
	if len(reorder) > 0 && reorder[0] {
		s.orderStr = value
	} else {
		s.orderStr = s.orderStr + value
	}
	return s
}

func (s *Orm) Count() int64 {
	return 0
}

func (s *Orm) Select(value interface{}) *Orm {
	switch value := value.(type) {
	case string:
		s.selectStr = value
	default:
		s.Error = errors.New("Can' understand the value of Select, Should be string")
	}
	return s
}

func (s *Orm) Save(value interface{}) *Orm {
	s.Model(value)
	if s.model.PrimaryKeyIsEmpty() {
		s.explain(value, "Create").create(value)
	} else {
		s.explain(value, "Update").update(value)
	}
	return s
}

func (s *Orm) Delete(value interface{}) *Orm {
	s.explain(value, "Delete").Exec()
	return s
}

func (s *Orm) Update(column string, value string) *Orm {
	return s
}

func (s *Orm) Updates(values map[string]string) *Orm {
	return s
}

func (s *Orm) Exec(sql ...string) *Orm {
	if len(sql) == 0 {
		s.SqlResult, s.Error = s.db.Exec(s.Sql, s.SqlVars...)
	} else {
		s.SqlResult, s.Error = s.db.Exec(sql[0])
	}
	return s
}

func (s *Orm) First(out interface{}) *Orm {
	s.explain(out, "Query").query(out)
	return s
}

func (s *Orm) Find(out interface{}) *Orm {
	s.explain(out, "Query").query(out)
	return s
}

func (s *Orm) Or(querystring interface{}, args ...interface{}) *Orm {
	return s
}

func (s *Orm) Not(querystring interface{}, args ...interface{}) *Orm {
	return s
}

func (s *Orm) CreateTable(value interface{}) *Orm {
	s.explain(value, "CreateTable").Exec()
	return s
}
