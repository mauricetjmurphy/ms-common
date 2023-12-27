package query_test

import (
	"context"
	"github.com/NBCUniversal/gvs-ms-common/db/dbmocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NBCUniversal/gvs-ms-common/db/entity"
	q "github.com/NBCUniversal/gvs-ms-common/db/query"
	c "github.com/NBCUniversal/gvs-ms-common/db/query/criteria"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	dummyTable = "dummy"
	oneToTable = "one_to"
)

type DummyTable struct {
	*entity.Base
	Name    string `gorm:"column:Name"`
	OneToID *uint  `gorm:"column:OneToId"`
	OneTo   *OneTo
	*entity.Audit
}

func (DummyTable) TableName() string { return dummyTable }

type OneTo struct {
	*entity.Base
	Status string `gorm:"column:Status"`
}

func (OneTo) TableName() string { return oneToTable }

var (
	allSelectDummyCols = []string{"Id", "Name", "OneToId"}
	allSelectCountCols = []string{"COUNT(DISTINCT(`dd`.`Id`))"}
)

type testCases struct {
	name      string
	fnMocks   func(sqlMock sqlmock.Sqlmock)
	query     func(db *gorm.DB) q.Querier
	wantTotal int64
	wantItems int64
}

func TestQuery_Find(t *testing.T) {
	cases := []struct {
		name    string
		fnMocks func(sqlMock sqlmock.Sqlmock)
		query   func(db *gorm.DB) q.Querier
		wantErr error
	}{
		{
			name: "Find_WithNoWhereClause",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd")
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Equal",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`ID` = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.Eq("dd.ID", 1))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_NotEqual",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`ID` <> ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.NotEq("dd.ID", 1))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_ISNULL",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`OneToID` IS NULL").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", nil))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.IsNil("dd.OneToID"))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_ISNOTNULL",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`OneToID` IS NOT NULL").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.IsNotNil("dd.OneToID"))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_In",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE \\(`dd`.`ID` = \\? AND `dd`.`Name` IN \\(\\?,\\?\\)\\)").
					WithArgs(1, "Test", "ABC").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.And(c.Eq("dd.ID", 1), c.In("dd.Name", "Test", "ABC")))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Empty_In",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE \\(`dd`.`ID` = \\? AND `dd`.`Name` IN \\(\\(NULL\\)\\)\\)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				var emptyStrs []string
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.And(c.Eq("dd.ID", 1), c.In("dd.Name", emptyStrs)))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_NotIn",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE \\(`dd`.`ID` = \\? AND `dd`.`Name` NOT IN \\(\\?,\\?\\)\\)").
					WithArgs(1, "Test", "ABC").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.And(c.Eq("dd.ID", 1), c.NotIn("dd.Name", "Test", "ABC")))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Empty_NotIn",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE \\(`dd`.`ID` = \\? AND `dd`.`Name` NOT IN \\(\\(NULL\\)\\)\\)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				var emptyStrs []string
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.And(c.Eq("dd.ID", 1), c.NotIn("dd.Name", emptyStrs)))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Like",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` LIKE ?").
					WithArgs("Test").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.Like("dd.Name", "Test"))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Wildcard",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` LIKE ?").
					WithArgs("%Test%").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.Wildcard("dd.Name", "Test"))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_MultipleCondition_And",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE \\(`dd`.`ID` = \\? AND `dd`.`Name` = \\?\\)").
					WithArgs(1, "Test").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.And(c.Eq("dd.ID", 1), c.Eq("dd.Name", "Test")))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_MultipleCondition_Or",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE \\(`dd`.`ID` = \\? OR `dd`.`Name` = \\?\\)").
					WithArgs(1, "Test").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.Or(c.Eq("dd.ID", 1), c.Eq("dd.Name", "Test")))
			},
			wantErr: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			//Given
			dbMock, sqlMock := dbmocks.NewSqlMock()
			tt.fnMocks(sqlMock)
			var results []*DummyTable

			//When
			err := tt.query(dbMock).Find(context.Background(), &results)

			//Then
			assert.Equal(t, tt.wantErr, err)
			// we make sure that all expectations were met
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuery_FindAll(t *testing.T) {
	//Given
	cases := []testCases{
		{
			name: "Select_Pagination_With_Distinct_True_Where_Clause_Wildcard_Search",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT COUNT\\(DISTINCT\\(`dd`.`Id`\\)\\) FROM dummy AS dd  WHERE `dd`.`Name` LIKE ?").
					WithArgs("%Test%").
					WillReturnRows(sqlmock.NewRows(allSelectCountCols).AddRow(1))

				sqlMock.ExpectQuery("SELECT DISTINCT dd.* FROM dummy AS dd  WHERE `dd`.`Name` LIKE \\? ORDER BY dd.Id LIMIT 10").
					WithArgs("%Test%").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))

			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					Distinct(true).
					SelectCount("dd.Id").
					From(dummyTable, "dd").
					Where([]c.Expr{
						c.Wildcard("dd.Name", "Test"),
					}...).
					Page(0, 10).
					Order("dd.Id")
			},
			wantTotal: 1,
			wantItems: 1,
		},
		{
			name: "Select_Pagination_Joins_Where_Clause_Wildcard_Search",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT COUNT\\(DISTINCT\\(`dd`.`Id`\\)\\) FROM dummy AS dd INNER JOIN `one_to` `ott` ON ott.ID = dd.OneToID WHERE `dd`.`Name` LIKE ?").
					WithArgs("%Test%").
					WillReturnRows(sqlmock.NewRows(allSelectCountCols).AddRow(1))

				sqlMock.ExpectQuery("SELECT DISTINCT dd.* FROM dummy AS dd INNER JOIN `one_to` `ott` ON ott.ID = dd.OneToID WHERE `dd`.`Name` LIKE \\? ORDER BY dd.Id LIMIT 10").
					WithArgs("%Test%").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))

			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					Distinct(true).
					SelectCount("dd.Id").
					From(dummyTable, "dd").
					InnerJoin(q.Join{Table: oneToTable, As: "ott", Cond: "ott.ID = dd.OneToID"}).
					Where([]c.Expr{
						c.Wildcard("dd.Name", "Test"),
					}...).
					Page(0, 10).
					Order("dd.Id")
			},
			wantTotal: 1,
			wantItems: 1,
		},
		{
			name: "Select_None_Distinct_Where_Clause_Wildcard_Search",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT COUNT\\(`dd`.*\\) FROM dummy AS dd  WHERE `dd`.`Name` LIKE ?").
					WithArgs("%test%").
					WillReturnRows(sqlmock.NewRows(allSelectCountCols).AddRow(1))

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd  WHERE `dd`.`Name` LIKE ?").
					WithArgs("%test%").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))

			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where([]c.Expr{
						c.Wildcard("dd.Name", "test"),
					}...)
			},
			wantTotal: 1,
			wantItems: 1,
		},
	}

	dbMock, sqlMock := dbmocks.NewSqlMock()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			//Given
			tt.fnMocks(sqlMock)
			var results []*DummyTable

			//When
			total, _ := tt.query(dbMock).FindAll(context.Background(), &results)

			//Then
			assert.Equal(t, tt.wantTotal, total)
			assert.Equal(t, tt.wantTotal, int64(len(results)))

			// we make sure that all expectations were met
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestQuery_Find_Utils(t *testing.T) {
	cases := []struct {
		name    string
		fnMocks func(sqlMock sqlmock.Sqlmock)
		query   func(db *gorm.DB) q.Querier
		wantErr error
	}{
		{
			name: "Find_WithWhereClause_Empty_InStr",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`ID` = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				var emptyStrs []string
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.And(c.Eq("dd.ID", 1), c.InStr("dd.Name", emptyStrs...)))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Value_InStr",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?\\)").
					WithArgs("Test").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InStr("dd.Name", "Test"))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Values_InStr",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?,\\?\\)").
					WithArgs("Test", "Test1").
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InStr("dd.Name", "Test", "Test1"))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Values_InUint",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?,\\?\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InUint("dd.Name", 1, 2))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Values_InUint32",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?,\\?\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InUint32("dd.Name", 1, 2))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Values_InUint64",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?,\\?\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InUint64("dd.Name", 1, 2))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Values_InInt",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?,\\?\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InInt("dd.Name", 1, 2))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Values_InInt32",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?,\\?\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InInt32("dd.Name", 1, 2))
			},
			wantErr: nil,
		},
		{
			name: "Find_WithWhereClause_Values_InInt64",
			fnMocks: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.MatchExpectationsInOrder(false)

				sqlMock.ExpectQuery("SELECT dd.* FROM dummy AS dd WHERE `dd`.`Name` IN \\(\\?,\\?\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows(allSelectDummyCols).AddRow(1, "Test", 1))
			},
			query: func(db *gorm.DB) q.Querier {
				return q.New(db).
					Select("dd.*").
					From(dummyTable, "dd").
					Where(c.InInt64("dd.Name", 1, 2))
			},
			wantErr: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			//Given
			dbMock, sqlMock := dbmocks.NewSqlMock()
			tt.fnMocks(sqlMock)
			var results []*DummyTable

			//When
			err := tt.query(dbMock).Find(context.Background(), &results)

			//Then
			assert.Equal(t, tt.wantErr, err)
			// we make sure that all expectations were met
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
