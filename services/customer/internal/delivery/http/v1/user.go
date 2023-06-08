package v1

import (
	"errors"
	"net/http"

	"github.com/forstes/besafe-go/customer/pkg/request"
	"github.com/forstes/besafe-go/customer/services/customer/internal/service"
)

type UserHandler struct {
	userService service.Users
}

func NewUserHandler(userService service.Users) *UserHandler {
	return &UserHandler{userService: userService}
}

type registerUserDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type loginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	var dto registerUserDTO

	if err := request.ReadJSON(w, r, &dto); err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	input := service.UserSignUpInput{
		Email:     dto.Email,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Phone:     dto.LastName,
	}

	err := h.userService.SignUp(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			request.BadRequestResponse(w, r, err)
			return
		case errors.Is(err, service.ErrDuplicate):
			request.RecordDuplicationResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var dto loginUserDTO

	if err := request.ReadJSON(w, r, &dto); err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	input := service.UserSignInInput{
		Email:    dto.Email,
		Password: dto.Password,
	}

	token, err := h.userService.SignIn(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWrongCredentials):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	request.WriteJSON(w, http.StatusOK, map[string]any{"token": token.PlainText}, nil)
}
