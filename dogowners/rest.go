package dogowners

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/almendar/golang-gorm-chi-postgres/shared"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Dog struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
}

type Owner struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Dogs []Dog  `json:"dogs"`
}

type AssignDogToOwnerRequest struct {
	DogName  string    `json:"dog_name" validate:"required"`
	DogBirth time.Time `json:"dog_birth" validate:"required"`
}

type HttpHandlers struct {
	R   *chi.Mux
	svc *Service
}

func NewHttpHandlers(service *Service) *HttpHandlers {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	handler := &HttpHandlers{svc: service, R: r}

	r.Get("/api/v1/owner", handler.ListOwners)
	r.Get("/api/v1/owner/{userID}", handler.GetOwner)
	r.Put("/api/v1/owner/{userID}/dog", handler.AssignDogToOwner)

	return handler
}

// ListOwners lists all owners
//
//	@Summary		List owners
//	@Description	get all owners
//	@Tags			owners
//	@Produce		json
//	@Success		200	{array}	dogowners.Owner
//	@Router			/owner [get]
func (h *HttpHandlers) ListOwners(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListOwners()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(users) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *HttpHandlers) GetOwner(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userIDUint, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetDogowner(uint(userIDUint))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(user) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *HttpHandlers) AssignDogToOwner(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userIDUint, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req AssignDogToOwnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if validationErrors, err := shared.RunValidation(req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to run validation: %s\n", err)
		return
	} else if validationErrors != nil {
		jsonRet := struct {
			Errors []shared.ValidationError `json:"errors"`
		}{validationErrors}

		w.WriteHeader(http.StatusBadRequest)
		if json.NewEncoder(w).Encode(jsonRet) != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if err := h.svc.SaveDog(uint(userIDUint), Dog{Name: req.DogName, Birthday: req.DogBirth}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to save dog: %s\n", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *HttpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.R.ServeHTTP(w, r)
}
