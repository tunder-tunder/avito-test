package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/tunder-tunder/avito/bd"
	"gitlab.com/tunder-tunder/avito/models"
	"net/http"
	"strconv"
)

var serviceIDKey = "serviceID"

func services(router chi.Router) {
	router.Get("/", getAllServices)
	router.Post("/", createService)
	router.Route("/{serviceId}", func(router chi.Router) {
		router.Use(ServiceContext)
		router.Get("/", getService)
		router.Delete("/", deleteService)
	})
}

func ServiceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceId := chi.URLParam(r, "serviceId")

		if serviceId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("service ID is required")))
			return
		}

		id, err := strconv.Atoi(serviceId)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid Service ID")))
		}

		ctx := context.WithValue(r.Context(), serviceIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createService(w http.ResponseWriter, r *http.Request) {
	item := &models.Service{}

	if err := render.Bind(r, item); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddService(item); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		//initBalance(serviceId)
		return
	}
	if err := render.Render(w, r, item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
func getAllServices(w http.ResponseWriter, r *http.Request) {
	items, err := dbInstance.GetAllServices()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, items); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getService(w http.ResponseWriter, r *http.Request) {
	serviceId := r.Context().Value(serviceIDKey).(int)
	item, err := dbInstance.GetServiceById(serviceId)
	if err != nil {
		if err == bd.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteService(w http.ResponseWriter, r *http.Request) {
	serviceId := r.Context().Value(serviceIDKey).(int)
	err := dbInstance.DeleteService(serviceId)
	if err != nil {
		if err == bd.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}
