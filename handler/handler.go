package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/tunder-tunder/avito/bd"
	_ "gitlab.com/tunder-tunder/avito/models"
	"net/http"
	_ "strconv"
)

var dbInstance bd.Database

func NewHandler(db bd.Database) http.Handler {
	router := chi.NewRouter()
	dbInstance = db
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/users", users)
	router.Route("/servcies", services)
	router.Route("/balances", balances)
	return router
}
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, ErrNotFound)
}
