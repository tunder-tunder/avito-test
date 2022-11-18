package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/tunder-tunder/avito/bd"
	"gitlab.com/tunder-tunder/avito/models"
	"log"
	"net/http"
	"strconv"
)

var userIdKey = "userId"
var balanceIdKey = "balanceId"

func balances(router chi.Router) {
	router.Get("/", getAllBalances)
	router.Route("/{userId}", func(router chi.Router) {
		router.Use(BalanceContext)
		//router.Post("/", initBalance)
		router.Post("/", addBalance)
		router.Get("/", getBalanceById)
		router.Post("/reserve", reserveBalance)
		router.Post("/pay", payBalance)
	})

}
func BalanceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")

		if userId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("user ID is required")))
			return
		}

		id_user, err := strconv.Atoi(userId)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid user ID")))
		}

		ctx_userId := context.WithValue(r.Context(), userIDKey, id_user)
		next.ServeHTTP(w, r.WithContext(ctx_userId))

	})
}

func addBalance(w http.ResponseWriter, r *http.Request) {
	topup := chi.URLParam(r, "top-up")
	if topup != "" {
		intTopUp, _ := strconv.Atoi(topup)
		userID := r.Context().Value(userIDKey).(int)
		item, err := dbInstance.GetBalanceById(userID)
		if err != nil {
			if err == bd.ErrNoMatch {
				render.Render(w, r, ErrNotFound)
			} else {
				newBalance := item.Total + intTopUp
				dbInstance.AddBalance(userID, newBalance)
				render.Render(w, r, ErrorRenderer(err))
			}
			return
		}
		if err := render.Render(w, r, &item); err != nil {
			render.Render(w, r, ServerErrorRenderer(err))
			return
		}
	}

}

func getAllBalances(w http.ResponseWriter, r *http.Request) {
	items, err := dbInstance.GetAllBalances()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, items); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getBalanceById(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	item, err := dbInstance.GetBalanceById(userID)
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

func reserveBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	balanceData := &models.Balance{}

	serviceId := chi.URLParam(r, "service-id")
	intServiceId, _ := strconv.Atoi(serviceId)
	price := chi.URLParam(r, "price")
	intPrice, _ := strconv.Atoi(price)
	orderNumber := chi.URLParam(r, "order-number")
	balance, _ := dbInstance.GetBalanceById(userId)
	total := balance.Total - intPrice

	if err := render.Bind(r, balanceData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	//userId int, serviceId int, orderNumber string, price int, total int
	item, err := dbInstance.ReserveBalance(userId, intServiceId, orderNumber, intPrice, total)
	if err != nil {
		if err == bd.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func payBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	balanceData := &models.Balance{}

	serviceId := chi.URLParam(r, "service-id")
	intServiceId, _ := strconv.Atoi(serviceId)

	orderNumber := chi.URLParam(r, "order-number")

	balance, _ := dbInstance.GetOrderByNumber(orderNumber)
	total := balance.Total
	reserve := balance.Reserve
	log.Println(reserve)

	if err := render.Bind(r, balanceData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	// userId int, serviceId int, orderNumber string, total int
	item, err := dbInstance.PayBalance(userId, intServiceId, orderNumber, total)
	if err != nil {
		if err == bd.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}

	if err := render.Render(w, r, item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
