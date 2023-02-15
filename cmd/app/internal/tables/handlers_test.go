package tables_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/database"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/tables"
)

// type v2Suite struct {
// 	db   *gorm.DB
// 	mock sqlmock.Sqlmock
// 	// table models.Table
// }

var (
	testDBURL = "postgresql://user:password@localhost:5432/database?sslmode=disable"
)

func TestCreate(t *testing.T) {
	db, err := database.ReadyStateDB(testDBURL)
	require.NoError(t, err)

	body := "{\"capacity\":10}"
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", "/tables", reader)
	require.NoError(t, err)

	r := tables.NewRepository(db)
	service := tables.NewService(r)
	c := tables.NewController(service)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.Create())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
