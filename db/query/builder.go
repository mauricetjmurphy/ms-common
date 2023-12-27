package query

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/mauricetjmurphy/ms-common/db/query/criteria"
)

// Join type
type Join struct {
	Table, As, Cond string
}

// Select provides the set of column to be used to retrieve row selected columns from one or more tables.
func (q *Query) Select(columns ...string) *Query {
	q.fields = columns
	return q
}

// Distinct provides the set of columns to be used to remove of duplicated from a result set.
func (q *Query) Distinct(distinct bool) *Query {
	q.distinct = distinct
	return q
}

// SelectCount provides the set of columns to be used to select count.
func (q *Query) SelectCount(columns ...string) *Query {
	q.countDistinct = columns
	return q
}

// From provides the source table name to be selected and possibly other clauses.
func (q *Query) From(table, alias string) *Query {
	q.from = table
	if len(alias) > 0 {
		q.from = fmt.Sprintf("%v AS %v", q.from, alias)
	}
	return q
}

// InnerJoin provides the inner innerJoin expression.
func (q *Query) InnerJoin(join Join) *Query {
	q.innerJoin = append(q.innerJoin, join)
	return q
}

// InnerJoinIf provides the inner join expression when matching the given conditional.
func (q *Query) InnerJoinIf(join Join, cond bool) *Query {
	if cond {
		return q.InnerJoin(join)
	}
	return q
}

// LeftJoin provides the left innerJoin expression.
func (q *Query) LeftJoin(left Join) *Query {
	q.leftJoin = append(q.leftJoin, left)
	return q
}

// LeftJoinIf provides the left innerJoin expression when matching the given conditional.
func (q *Query) LeftJoinIf(left Join, cond bool) *Query {
	if cond {
		return q.LeftJoin(left)
	}
	return q
}

// Where provides array of expressions with one more conditions
// that evaluate to true for each row to be selected.
func (q *Query) Where(exps ...criteria.Expr) *Query {
	var where []criteria.Expr
	for _, expr := range exps {
		if expr != nil {
			where = append(where, expr)
		}
	}
	q.where = where
	return q
}

// Order provides the ordering rows in a result set.
func (q *Query) Order(raw string) *Query {
	if len(raw) > 0 {
		q.order = raw
	}
	return q
}

// GroupBy provides the group by rows in a result set.
func (q *Query) GroupBy(raw string) *Query {
	if len(raw) > 0 {
		q.group = raw
	}
	return q
}

// Offset provides the page size to be returned a result set.
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// Limit provides the page size to be returned a result set.
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// Page provides the basic to be returned a pagination result set on given the offset and limit.
func (q *Query) Page(offset, limit int64) *Query {
	q.Limit(int(limit))
	q.Offset(int(offset) * q.limit)
	return q
}

// Preloads loads relation entity models to be allowed eager reloading.
// See https://gorm.io/docs/preload.html for more usage.
func (q *Query) Preloads(models ...string) *Query {
	q.preloads = models
	return q
}

// ReloadFunc provides way to load reloation entity models with conditions
// See https://gorm.io/docs/preload.html for more usage.
func (q *Query) ReloadFunc(model string, criteria criteria.Expr) *Query {
	if criteria == nil {
		return q
	}
	if _, ok := q.preloadFunc[model]; !ok {
		q.preloadFunc[model] = criteria
	}
	return q
}

// ReloadConds provides way to load reloation entity models with a map conditions.
// See https://gorm.io/docs/preload.html for more usage.
func (q *Query) ReloadConds(conditions map[string]criteria.Expr) *Query {
	if len(conditions) == 0 {
		return q
	}
	for k, v := range conditions {
		q.ReloadFunc(k, v)
	}
	return q
}

func buildJoin(joins []Join, joinType clause.JoinType) []clause.Join {
	var joinClauses []clause.Join
	for _, j := range joins {
		joinClauses = append(joinClauses, clause.Join{
			Type:  joinType,
			Table: clause.Table{Name: j.Table, Alias: j.As},
			ON:    clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: j.Cond}}},
		})
	}
	return joinClauses
}
