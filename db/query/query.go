package query

import (
	"context"

	"github.com/mauricetjmurphy/ms-common/db/query/criteria"

	"golang.org/x/sync/errgroup"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Querier interface presents the propagation query functions to database.
//
//go:generate mockery --output querymocks --outpkg querymocks --name Querier
type Querier interface {
	Find(context.Context, interface{}) error
	FindAll(context.Context, interface{}) (int64, error)
}

type Query struct {
	db            *gorm.DB
	model         interface{}
	fields        []string
	distinct      bool
	countDistinct []string
	from          string
	innerJoin     []Join
	leftJoin      []Join
	where         []criteria.Expr
	order         string
	group         string
	offset        int
	limit         int
	preloads      []string
	preloadFunc   map[string]criteria.Expr
}

// New create the query instance on given client DB.
func New(db *gorm.DB) *Query {
	return &Query{
		db:          db,
		preloadFunc: make(map[string]criteria.Expr),
	}
}

func (q *Query) FindAll(ctx context.Context, value interface{}) (int64, error) {
	q.context(ctx)

	var (
		eg    errgroup.Group
		count int64
	)

	eg.Go(func() error {
		// Count the total records matching given criteria
		return q.getTotalRecords(value, &count)
	})

	eg.Go(func() error {
		// Find paged records matching given criteria
		result := q.buildClauses(value).Find(value)
		return result.Error
	})

	if err := eg.Wait(); err != nil {
		return 0, err
	}

	return count, nil
}

func (q *Query) Find(ctx context.Context, value interface{}) error {
	q.context(ctx)
	result := q.buildClauses(value).Find(value)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (q *Query) Clause() *gorm.DB {
	tx := q.db.Model(q.model)

	tx.Select(q.fields)
	if q.distinct {
		tx.Distinct()
	}

	tx.Table(q.from)
	tx.Clauses(
		clause.From{
			Joins: q.withJoin(),
		},
	)

	if len(q.where) > 0 {
		tx.Clauses(clause.Where{
			Exprs: q.where,
		})
	}

	if q.limit > 0 {
		tx.Limit(q.limit)
		tx.Offset(q.offset)
	}

	if len(q.order) > 0 {
		tx.Order(q.order)
	}

	if len(q.group) > 0 {
		tx.Group(q.group)
	}

	if len(q.preloadFunc) > 0 {
		for k, v := range q.preloadFunc {
			tx.Preload(k, v)
		}
	}

	if len(q.preloads) > 0 {
		for _, m := range q.preloads {
			tx.Preload(m)
		}
	}

	return tx
}

func (q *Query) context(ctx context.Context) {
	if ctx != nil {
		q.db.WithContext(ctx)
	}
}

func (q *Query) withModel(model interface{}) *Query {
	q.model = model
	return q
}

func (q *Query) buildClauses(model interface{}) *gorm.DB {
	return q.withModel(model).Clause()
}

func (q *Query) withJoin() []clause.Join {
	inner := buildJoin(q.innerJoin, clause.InnerJoin)
	left := buildJoin(q.leftJoin, clause.LeftJoin)
	return append(inner, left...)
}

func (q *Query) getTotalRecords(value interface{}, count *int64) error {
	cli := q.buildClauses(value)
	if q.distinct && len(q.countDistinct) > 0 {
		cli.Distinct(q.countDistinct)
	}
	return cli.Count(count).Error
}
