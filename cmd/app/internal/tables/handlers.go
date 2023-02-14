package tables

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
)

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (c Controller) ListTables() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tableInfo, err := c.service.GetTables(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(tableInfo)

	}
}

func (c Controller) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		table := models.Table{}
		err = json.Unmarshal(body, &table)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		err = c.service.CreateTable(r.Context(), table)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// func (c Controller) GetTables() http.HandlerFunc {
// 	// var tab table
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// tableInfo := p.db.QueryRow("`SELECT * FROM tables WHERE id = 1;`")
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 		// 	return
// 		// }
// 		// err := tableInfo.Scan(&tab.id, &tab.capacity)
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 		// 	return
// 		// }

// 		for _, table := range models.Tables {
// 			// str := fmt.Sprintf("id: %s, capacity: %d", table.ID, table.Capacity)
// 			fmt.Fprintf(w, table.ID)
// 			fmt.Fprintf(w, "\n")
// 			fmt.Fprint(w, table.Capacity)
// 			fmt.Fprintf(w, "\n")
// 			fmt.Printf(table.ID)
// 			fmt.Printf("\n")
// 			fmt.Print(table.Capacity)
// 			fmt.Printf("\n")
// 		}

// 		fmt.Fprintf(w, "omg get worked maybe\n")
// 	}
// }

// func (c Controller) getTable(table_id string) (int, error) {
// 	for _, table := range models.Tables {
// 		if table.ID == table_id {
// 			return table.Capacity, nil
// 		}
// 	}

// 	return 0, errors.New("sorry this table doesn't exist, try again")
// }

// func (c Controller) checkTableCapacity(capacity int, partySize int) bool {
// 	return capacity >= partySize
// }

// // Create returns an handler func which creates a new table
// func (c Controller) Create() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// ctx := r.Context()
// 		// body, err := io.ReadAll(r.Body)
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 		// 	return
// 		// }
// 		// table := models.Table{}
// 		// err = json.Unmarshal(body, &table)
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusForbidden)
// 		// 	return
// 		// }
// 		c.service.CreateTable()

// 		fmt.Fprintf(w, "added table")

// 		w.WriteHeader(http.StatusOK)
// 	}
// }

// func (c Controller) CalculateEmptySeats() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		totalEmptySeats := 0
// 		for _, table := range models.Tables {
// 			totalEmptySeats += table.SeatsEmpty
// 		}
// 		// totalEmptySeats, err := p.db.Exec(`SUM(empty_seats) FROM tables;`)
// 		// if err != nil {
// 		// 	w.WriteHeader(http.StatusInternalServerError)
// 		// }
// 		// fmt.Fprint(w, err.Error(), totalEmptySeats)
// 		fmt.Fprint(w, totalEmptySeats)
// 		w.WriteHeader(http.StatusOK)
// 	}
// }

func (c Controller) HandlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong\n")
	}
}
