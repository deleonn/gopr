package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joeldeleon/pr_description_generator/internal/models"
)

type API struct {
	router *mux.Router
}

func NewAPI() *API {
	r := mux.NewRouter()
	api := &API{router: r}

	// Register routes
	api.registerRoutes()

	return api
}

func (a *API) registerRoutes() {
	a.router.HandleFunc("/generate-pr", a.generatePRDescription).Methods("POST")
}

func (a *API) generatePRDescription(w http.ResponseWriter, r *http.Request) {
	var req models.PRRequest
	if err := parseJSON(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Generate PR description using LLM
	description, err := generatePRDescription(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate PR description")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"description": description})
}

func (a *API) Start(address string) error {
	return http.ListenAndServe(address, a.router)
}

