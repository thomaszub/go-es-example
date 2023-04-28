package api

import (
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/thomaszub/go-es-example/domain"
)

type AccountController struct {
	service *domain.AccountService
}

func NewAccountController(service *domain.AccountService) AccountController {
	return AccountController{
		service: service,
	}
}

func (c *AccountController) RegisterOn(baseRoute *echo.Group) {
	baseRoute.GET("", c.GetAccounts)
	baseRoute.POST("", c.CreateAccount)
	baseRoute.GET("/:id", c.GetAccount)
	baseRoute.DELETE("/:id", c.DeleteAccount)
	baseRoute.POST("/:id/deposit", c.Deposit)
	baseRoute.POST("/:id/withdraw", c.Withdraw)
	baseRoute.PUT("/:id/limit", c.SetLimit)
}

type getAccountsResponse struct {
	AccountIds []gocql.UUID `json:"accountIds"`
}

func (c *AccountController) GetAccounts(ctx echo.Context) error {
	ids, err := c.service.GetAllAccountIds()
	if err != nil {
		return domainError(err)
	}
	return ctx.JSON(http.StatusOK, getAccountsResponse{AccountIds: ids})
}

type newAccountResponse struct {
	AccountId gocql.UUID `json:"accountId"`
}

func (c *AccountController) CreateAccount(ctx echo.Context) error {
	acc, err := c.service.CreateNewAccount()
	if err != nil {
		return domainError(err)
	}
	return ctx.JSON(http.StatusCreated, newAccountResponse{AccountId: acc.AccountId()})
}

type getAccountResponse struct {
	AccountId gocql.UUID `json:"accountId"`
	Limit     float64    `json:"limit"`
	Balance   float64    `json:"balance"`
}

func (c *AccountController) GetAccount(ctx echo.Context) error {
	id, err := getId(ctx)
	if err != nil {
		return err
	}
	acc, err := c.service.GetAccount(id)
	if err != nil {
		return domainError(err)
	}
	return ctx.JSON(http.StatusOK, getAccountResponse{
		AccountId: id,
		Limit:     acc.Limit(),
		Balance:   acc.Balance(),
	})
}

type depositRequest struct {
	Amount float64 `json:"amount"`
}

func (c *AccountController) Deposit(ctx echo.Context) error {
	id, err := getId(ctx)
	if err != nil {
		return err
	}
	body := depositRequest{}
	if err := ctx.Bind(&body); err != nil {
		return badRequest(err, err.Error())
	}
	acc, err := c.service.GetAccount(id)
	if err != nil {
		return domainError(err)
	}
	if err := acc.Deposit(body.Amount); err != nil {
		return domainError(err)
	}
	return ctx.NoContent(http.StatusAccepted)
}

type withdrawRequest struct {
	Amount float64 `json:"amount"`
}

func (c *AccountController) Withdraw(ctx echo.Context) error {
	id, err := getId(ctx)
	if err != nil {
		return err
	}
	body := withdrawRequest{}
	if err := ctx.Bind(&body); err != nil {
		return badRequest(err, err.Error())
	}
	acc, err := c.service.GetAccount(id)
	if err != nil {
		return domainError(err)
	}
	if err := acc.Withdraw(body.Amount); err != nil {
		return domainError(err)
	}
	return ctx.NoContent(http.StatusAccepted)
}

type setLimitRequest struct {
	Limit float64 `json:"limit"`
}

func (c *AccountController) SetLimit(ctx echo.Context) error {
	id, err := getId(ctx)
	if err != nil {
		return err
	}
	body := setLimitRequest{}
	if err := ctx.Bind(&body); err != nil {
		return badRequest(err, err.Error())
	}
	acc, err := c.service.GetAccount(id)
	if err != nil {
		return domainError(err)
	}
	if err := acc.SetNewLimit(body.Limit); err != nil {
		return domainError(err)
	}
	return ctx.NoContent(http.StatusAccepted)
}

func getId(ctx echo.Context) (gocql.UUID, error) {
	idString := ctx.Param("id")
	id, err := gocql.ParseUUID(idString)
	if err != nil {
		return gocql.UUID{}, badRequest(err, fmt.Sprintf("%s is not a valid id", idString))
	}
	return id, nil
}

func badRequest(err error, message string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     http.StatusBadRequest,
		Message:  message,
		Internal: err,
	}
}

func domainError(err error) *echo.HTTPError {
	code := http.StatusInternalServerError
	switch err.(type) {
	case *domain.AccountNotFoundError:
		code = http.StatusNotFound
	case *domain.DomainError:
		code = http.StatusBadRequest
	}
	return &echo.HTTPError{
		Code:     code,
		Message:  err.Error(),
		Internal: err,
	}
}

func (c *AccountController) DeleteAccount(ctx echo.Context) error {
	id, err := getId(ctx)
	if err != nil {
		return err
	}
	acc, err := c.service.GetAccount(id)
	if err != nil {
		return domainError(err)
	}
	err = acc.Delete()
	if err != nil {
		return domainError(err)
	}
	return ctx.NoContent(http.StatusNoContent)
}
