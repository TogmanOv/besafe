package v1

import (
	"net/http"

	"github.com/forstes/besafe-go/customer/services/customer/internal/service"
	"github.com/julienschmidt/httprouter"
)

type router struct {
	user UserHandler
}

func NewRouter(userService service.Users) *router {
	return &router{user: *NewUserHandler(userService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/user/register", r.user.Register)
	router.HandlerFunc(http.MethodPost, "/v1/user/login", r.user.Login)

	return router
}
