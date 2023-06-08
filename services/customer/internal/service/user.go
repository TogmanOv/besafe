package service

import (
	"context"
	"errors"
	"time"

	"github.com/forstes/besafe-go/customer/pkg/hash"
	token "github.com/forstes/besafe-go/customer/pkg/token"
	"github.com/forstes/besafe-go/customer/pkg/validator"
	"github.com/forstes/besafe-go/customer/services/customer/internal/domain"
	"github.com/forstes/besafe-go/customer/services/customer/internal/repository"
)

type UserSignUpInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Token struct {
	PlainText string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Token, error)
}

type userService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager token.TokenManager
}

func NewUserService(repo repository.Users, hasher hash.PasswordHasher, tokenManager token.TokenManager) *userService {
	return &userService{repo: repo, hasher: hasher, tokenManager: tokenManager}
}

func (s *userService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:        input.Email,
		PasswordHash: passwordHash,
		Details: domain.UserDetails{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Phone:     input.Phone,
		},
	}

	v := validator.New()
	validatePassword(v, input.Password)
	validateUser(v, &user)
	if !v.Valid() {
		// TODO More detailed error message
		return ErrFailedValidation
	}

	err = s.repo.Create(ctx, &user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicate):
			return ErrDuplicate
		default:
			return err
		}
	}
	return nil
}

func (s *userService) SignIn(ctx context.Context, input UserSignInInput) (Token, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Token{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return Token{}, ErrWrongCredentials
		default:
			return Token{}, err
		}
	}

	token, err := s.tokenManager.NewToken(user.ID, time.Duration(12)*time.Hour)
	return Token{PlainText: token}, err
}

func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func validatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func validateDetails(v *validator.Validator, ud *domain.UserDetails) {
	v.Check(ud.FirstName != "", "first name", "first name must be provided")
	v.Check(len(ud.FirstName) <= 100, "first name", "first name must not be more than 100 bytes long")

	v.Check(ud.LastName != "", "last name", "last name must be provided")
	v.Check(len(ud.LastName) <= 100, "last name", "last name must not be more than 100 bytes long")

	v.Check(ud.Phone != "", "phone", "phone must be provided")
}

func validateUser(v *validator.Validator, user *domain.User) {
	validateEmail(v, user.Email)
	validateDetails(v, &user.Details)
}
