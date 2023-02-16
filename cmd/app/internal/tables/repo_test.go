package tables_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/tables"

	// "github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	ctx  context.Context

	repository tables.Repository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)
	s.repository = tables.NewRepository(s.DB)
	s.ctx = context.Background()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_repository_Get() {
	var (
		id         = 1
		capacity   = 10
		seatsEmpty = capacity
	)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tables"`)).
		WillReturnRows(sqlmock.NewRows([]string{"ID", "Capacity", "SeatsEmpty"}).
			AddRow(id, capacity, seatsEmpty))
	s.mock.ExpectCommit()

	res, err := s.repository.GetTables(s.ctx)
	require.NoError(s.T(), err)

	realTable := []*models.Table{}
	err = json.Unmarshal(res, &realTable)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), models.Table{ID: id, Capacity: capacity, SeatsEmpty: seatsEmpty}, *realTable[0])
}

func (s *Suite) Test_repo_Create() {
	var (
		id         = 1
		capacity   = 45
		seatsEmpty = capacity
	)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tables" ("capacity","seats_empty") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(capacity, seatsEmpty).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectCommit()

	t := models.Table{}
	t.Capacity = 45
	t.SeatsEmpty = t.Capacity
	err := s.repository.CreateTable(s.ctx, t)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_repo_get_last_table_made() {
	var (
		id         = 1
		capacity   = 45
		seatsEmpty = capacity
	)
	// s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tables" ORDER BY "tables"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"ID", "Capacity", "SeatsEmpty"}).
			AddRow(id, capacity, seatsEmpty))
	// s.mock.ExpectCommit()

	_, err := s.repository.GetLastTableMade()
	require.NoError(s.T(), err)

}

func (s *Suite) Test_repo_empty_seats() {
	var (
		seatsEmpty = 10
	)
	// s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT SUM(seats_empty) FROM "tables"`)).
		WillReturnRows(sqlmock.NewRows([]string{"SeatsEmpty"}).
			AddRow(seatsEmpty))
	// s.mock.ExpectCommit()

	emptySeats, err := s.repository.CountNumberOfEmptySeats()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 10, emptySeats)
}

// func (s *Suite) Test_repo_edit_empty_seats_after_guest_leaves() {
// 	var (
// 		capacity = 10
// 		tableID = 1
// 		db = s.DB
// 	)
// 	// s.mock.ExpectBegin()
// 	s.mock.ExpectQuery(regexp.QuoteMeta(
// 		`UPDATE "tables" SET "seats_empty"=$1 WHERE id = $2"`)).
// 		WithArgs(capacity, tableID)
// 		// WillReturnRows(sqlmock.NewRows([]string{"SeatsEmpty"}).
// 		// 	AddRow(seatsEmpty))
// 	// s.mock.ExpectCommit()

// 	err := tables.EditEmptySeatsAfterGuestsLeave(capacity, tableID, db)
// 	require.NoError(s.T(), err)
// }
