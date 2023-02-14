package guests

import (
	"encoding/json"
	"io"

	// "io"
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

func (c Controller) AddGuestToGuestlist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		guest := models.Guest{}
		err = json.Unmarshal(body, &guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.service.AddGuestToGuestList(r.Context(), guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (c Controller) GetGuestsOnGuestList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestsOnList, err := c.service.GetGuestsOnGuestList(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(guestsOnList)
	}
}

func (c Controller) EditGuestsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// this handler should get called when guests arrive
		// we will ammend guest list if number of accompyaning guests+1 still fits in the capacity &
		// we will update the time of arrival as well

		// if too many people for table capacity, they get turned away and time of arrival doesn't get set 

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		guest := models.Guest{}
		err = json.Unmarshal(body, &guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		err = c.service.EditGuestsList(r.Context(), guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// func (c Controller) DeleteGuestFromGuestList() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }

