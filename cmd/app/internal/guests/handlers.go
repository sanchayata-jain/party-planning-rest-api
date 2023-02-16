package guests

import (
	"encoding/json"
	"io"

	"net/http"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
	"github.com/go-chi/chi"
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

		name := chi.URLParam(r, "name")

		guest := models.Guest{}
		err = json.Unmarshal(body, &guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		guest.Name = name

		if guest.AccompanyingGuests < 0 {
			w.Write([]byte("You can't have negative accompanying guests"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		responseName, err := c.service.AddGuestToGuestList(r.Context(), guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte(responseName))
		w.WriteHeader(http.StatusOK)
	}
}

func (c Controller) GetGuestsOnGuestList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestsNamesOnList, err := c.service.GetGuestsOnGuestList(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(guestsNamesOnList)
	}
}
// EditGuestList gets called when a guest wants to check into the party
func (c Controller) EditGuestsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := chi.URLParam(r, "name")

		guest := models.Guest{}
		err = json.Unmarshal(body, &guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		guest.Name = name

		if guest.AccompanyingGuests < 0 {
			w.Write([]byte("You can't have negative accompanying guests"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = c.service.EditGuestsList(r.Context(), guest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func (c Controller) DeleteGuestFromList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		err := c.service.DeleteGuestFromList(r.Context(), name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (c Controller) GetArrivedGuests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arrivedGuests, err := c.service.GetArrivedGuests(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.Write(arrivedGuests)
		w.WriteHeader(http.StatusAccepted)
	}
}
