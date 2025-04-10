package handler

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/kstsm/pvz-service/internal/auth"
	"github.com/kstsm/pvz-service/internal/middleware"
	"github.com/kstsm/pvz-service/internal/service"
	"net/http"
)

type HandlerI interface {
	NewRouter() http.Handler
	dummyLoginHandler(w http.ResponseWriter, r *http.Request)
	createPVZHandler(w http.ResponseWriter, r *http.Request)
	createReceptionHandler(w http.ResponseWriter, r *http.Request)
	addProductToReceptionHandler(w http.ResponseWriter, r *http.Request)
	deleteLastProductHandler(w http.ResponseWriter, r *http.Request)
	closeLastReceptionHandler(w http.ResponseWriter, r *http.Request)
	getListPVZ(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	ctx     context.Context
	service service.ServiceI
}

func NewHandler(ctx context.Context, svc service.ServiceI) HandlerI {
	return &Handler{
		ctx:     ctx,
		service: svc,
	}
}

func (h Handler) NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/dummyLogin", h.dummyLoginHandler)

	r.With(middleware.RoleMiddleware(auth.ValidateToken, "moderator")).Post("/pvz", h.createPVZHandler)
	r.With(middleware.RoleMiddleware(auth.ValidateToken, "pvz_employee")).Post("/receptions", h.createReceptionHandler)
	r.With(middleware.RoleMiddleware(auth.ValidateToken, "pvz_employee")).Post("/products", h.addProductToReceptionHandler)
	r.With(middleware.RoleMiddleware(auth.ValidateToken, "pvz_employee")).Post("/pvz/{pvzId}/delete_last_product", h.deleteLastProductHandler)
	r.With(middleware.RoleMiddleware(auth.ValidateToken, "pvz_employee")).Post("/pvz/{pvzId}/close_last_reception", h.closeLastReceptionHandler)
	r.With(middleware.RoleMiddleware(auth.ValidateToken, "pvz_employee", "moderator")).Get("/pvz", h.getListPVZ)
	return r
}
