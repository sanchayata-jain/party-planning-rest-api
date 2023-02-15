package tables

import (
	"encoding/json"
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
		w.WriteHeader(http.StatusOK)
	}
}

func (c Controller) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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

func (c Controller) SumEmptySeats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numOfEmptySeats, err := c.service.CountNumberOfEmptySeats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write([]byte(numOfEmptySeats))
		w.WriteHeader(http.StatusOK)
	}
}
