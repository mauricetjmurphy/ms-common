package criteria

import (
	"time"

	"gorm.io/gorm/clause"
)

// Expr alias gorm clause expression
type Expr = clause.Expression

func And(exprs ...Expr) Expr {
	return clause.And(ignoreNilExprs(exprs)...)
}

func Or(exprs ...Expr) Expr {
	return clause.Or(ignoreNilExprs(exprs)...)
}

func Eq(column string, value interface{}) Expr {
	return clause.Eq{Column: column, Value: value}
}

func NotEq(column string, value interface{}) Expr {
	return clause.Neq{Column: column, Value: value}
}

func IsNil(column string) Expr {
	return Eq(column, nil)
}

func IsNotNil(column string) Expr {
	return NotEq(column, nil)
}

func In(column string, values ...interface{}) Expr {
	return Eq(column, values)
}

func NotIn(column string, values ...interface{}) Expr {
	return NotEq(column, values)
}

func Like(column string, value string) Expr {
	if len(value) > 0 {
		return clause.Like{Column: column, Value: value}
	}
	return nil
}

// Wildcard clauses is the search wildcard expression value. Ex: LIKE '%keyword%'
func Wildcard(column, keyword string) Expr {
	if len(keyword) > 0 {
		return clause.Like{Column: column, Value: "%" + keyword + "%"}
	}
	return nil
}

func EqUintIfPresent(column string, value uint) Expr {
	if value > 0 {
		return Eq(column, value)
	}
	return nil
}

// SQLTimeExp clauses is SQL express with the datetime
func SQLTimeExp(sql string, value *time.Time, layout string) Expr {
	if value != nil && len(layout) > 0 {
		var format = value.Format(layout)
		return clause.Expr{SQL: sql, Vars: []interface{}{format}}
	}
	return nil
}

func Order(column, order string) (raw string) {
	if len(column) > 0 {
		raw = column
		if len(order) > 0 {
			raw = raw + " " + order
		}
	}
	return
}

func IsFalse(column string) Expr {
	return Eq(column, false)
}

func IsTrue(column string) Expr {
	return Eq(column, true)
}

func AlwaysTrue() Expr {
	return clause.NamedExpr{SQL: "1 = 1"}
}

func ignoreNilExprs(exprs []Expr) []Expr {
	var ignoreNilExpr []Expr
	for _, e := range exprs {
		if e != nil {
			ignoreNilExpr = append(ignoreNilExpr, e)
		}
	}
	return ignoreNilExpr
}
