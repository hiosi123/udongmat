package handlers

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/hiosi123/udongmat/auth"
	"github.com/hiosi123/udongmat/storage"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	Storage        *storage.AccountStorage
	SessionManager *auth.SessionManager
}

func NewAccountHandler(storage *storage.AccountStorage, sessionManager *auth.SessionManager) *AccountHandler {
	return &AccountHandler{Storage: storage, SessionManager: sessionManager}
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

// SignOut godoc
//
// @Summary Sign out
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /accounts/sign-out [post]
func (u *AccountHandler) SignOut(c *fiber.Ctx) error {
	// get the session from the authorization header
	sessionHeader := c.Get("Authorization")

	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.JSON(fiber.Map{"error": "invalid session header"})
	}

	// get the session id
	sessionId := sessionHeader[7:]

	// delete the session
	err := u.SessionManager.SignOut(sessionId)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

// GetInfo godoc
// @Summary Get the Accounts's info
// @Tags accounts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} auth.AccountSession
// @Router /accounts/me [get]
func (u *AccountHandler) GetInfo(c *fiber.Ctx) error {
	// get the session from the authorization header
	sessionHeader := c.Get("Authorization")

	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.JSON(fiber.Map{"error": "invalid session header"})
	}

	// get the session id
	sessionId := sessionHeader[7:]

	// get the account data from the session
	account, err := u.SessionManager.GetSession(sessionId)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(account)
}

type signInRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// SignInaccount godoc
// @Summary Sign in a account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body signInRequestBody true "The account's email and password"
// @Success 200 {object} SuccessResponse
// @Header 200 {string} Authorization "contains the session id in bearer format"
// @Router /accounts/sign-in [post]
func (u *AccountHandler) SignIn(c *fiber.Ctx) error {
	var account signInRequestBody

	err := c.BodyParser(&account)
	if err != nil {
		return err
	}

	fmt.Println(account)

	// validate the account struct
	validate := validator.New()
	err = validate.Struct(account)
	if err != nil {
		return err
	}

	// sign the account in
	sessionId, err := u.SessionManager.SignIn(account.Email, account.Password)
	if err != nil {
		return err
	}

	// set the session id as a header
	c.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	return c.JSON(fiber.Map{"success": true})
}

type AccountRequestBody struct {
	Name     string `json:"name" validate:"required,name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type signUpSuccessResponse struct {
	Id int `json:"id"`
}

// SignUpaccount godoc
// @Summary Sign up a account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body accountRequestBody true "The account's name, email, and password"
// @Success 200 {object} signUpSuccessResponse
// @Header 200 {string} Authorization "contains the session id in bearer format"
// @Router /accounts/sign-up [post]
func (u *AccountHandler) SignUp(c *fiber.Ctx) error {
	// get the info from the request body
	var account AccountRequestBody

	err := c.BodyParser(&account)
	if err != nil {
		return err
	}

	// validate the account struct
	validate := validator.New()
	err = validate.Struct(&account)
	if err != nil {
		return err
	}

	// hash the passowrd
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// create the account
	id, err := u.Storage.CreateNewAccount(storage.NewAccount{
		Name:     account.Name,
		Email:    account.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		return err
	}

	// generate the account's session
	sessionId, err := u.SessionManager.GenerateSession(auth.AccountSession{
		Id:    id,
		Name:  account.Name,
		Email: account.Email,
	})
	if err != nil {
		return err
	}

	// set the session id as a header
	c.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	resp := signUpSuccessResponse{
		Id: id,
	}

	return c.JSON(resp)
}
