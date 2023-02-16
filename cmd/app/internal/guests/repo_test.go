package guests_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/guests"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	ctx  context.Context

	repository guests.Repository
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
	s.repository = guests.NewRepository(s.DB)
	s.ctx = context.Background()
}

func (s *Suite) AfterTest(_, _ string) {
	s.mock.MatchExpectationsInOrder(false)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetGuestsOnGuestList() {
	var (
		name               = "Sanchayata"
		table              = 1
		accompanyingGuests = 4
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "guests"`)).
		WillReturnRows(sqlmock.NewRows([]string{"Name", "Table", "AccompanyingGuests"}).
			AddRow(name, table, accompanyingGuests))

	res, err := s.repository.GetGuestsOnGuestList(s.ctx)
	require.NoError(s.T(), err)

	guest := []*models.Guest{}
	err = json.Unmarshal(res, &guest)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), models.Guest{Name: name, Table: table, AccompanyingGuests: accompanyingGuests}, *guest[0])
}

func (s *Suite) TestGetGuestTableID() {
	var (
		name  = "sanchayata"
		table = 5
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "guests"`)).
		WillReturnRows(sqlmock.NewRows([]string{"Table"}).
			AddRow(table))

	res, err := guests.GetGuestTableID(s.ctx, s.DB, name)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), table, res)
}

func (s *Suite) TestAddGuestToGuestList() {
	var (
		name               = "San"
		table              = 1
		accompanyingGuests = 4
		arrived            = false
	)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "guests" ("name","table","accompanying_guests","arrived") VALUES ($1,$2,$3,$4)`)).
		WithArgs(name, table, accompanyingGuests, arrived).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	g := models.Guest{}
	g.Name = name
	g.Table = table
	g.AccompanyingGuests = accompanyingGuests
	g.Arrived = false
	responseName, err := s.repository.AddGuestToGuestlist(s.ctx, g)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), name, responseName)
}

func (s *Suite) TestEditGuestList() {
	var (
		name               = "Jason"
		timeArrived        = time.Now()
		accompanyingGuests = 2
		arrived            = true
	)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "guests" SET "accompanying_guests"=$1,"time_arrived"=$2,"arrived"=$3 WHERE name = $4`)).
		WithArgs(accompanyingGuests, timeArrived, arrived, name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	g := models.Guest{}
	g.Name = name
	g.AccompanyingGuests = accompanyingGuests
	g.Arrived = true
	err := s.repository.EditGuestList(timeArrived, g)
	require.NoError(s.T(), err)
}

func (s *Suite) TestDeleteGuest() {
	var (
		name = "Jason"
	)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "guests" WHERE name = $1`)).
		WithArgs(name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.DeleteGuest(name)
	require.NoError(s.T(), err)
}

func (s *Suite) TestGetArrivedGuests() {
	var (
		name               = "Jason"
		accompanyingGuests = 4
		timeArrived        = time.Date(2023, time.February, 16, 10, 9, 51, 595434000, time.UTC)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "guests" WHERE arrived = $1`)).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"Name", "AccompanyingGuests", "TimeArrived"}).
			AddRow(name, accompanyingGuests, timeArrived))

	res, err := s.repository.GetArrivedGuests(s.ctx)
	require.NoError(s.T(), err)

	guest := []*models.Guest{}
	err = json.Unmarshal(res, &guest)
	require.NoError(s.T(), err)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), models.Guest{Name: name, AccompanyingGuests: accompanyingGuests, TimeArrived: timeArrived}, *guest[0])
}
