package domain

import (
	"github.com/gocql/gocql"
)

type Account struct {
	repo      AccountEventRepository
	accountId gocql.UUID
	deleted   bool
	limit     float64
	balance   float64
}

func (a *Account) Deleted() bool {
	return a.deleted
}

func (a *Account) AccountId() gocql.UUID {
	return a.accountId
}

func (a *Account) Balance() float64 {
	return a.balance
}

func (a *Account) Limit() float64 {
	return a.limit
}

func (a *Account) SetNewLimit(limit float64) error {
	if limit > 0 {
		return NewDomainError("new limit %f can not be positive", limit)
	}
	if a.balance < limit {
		return NewDomainError("new limit %f can not be set as balance %f would be below limit", limit, a.balance)
	}
	e := LimitSetEvent{
		AccountId: a.accountId,
		EventId:   gocql.TimeUUID(),
		Limit:     limit,
	}
	if err := a.repo.Write(e); err != nil {
		return err
	}
	a.limit = limit
	return nil
}

func (a *Account) Deposit(amount float64) error {
	if amount < 0 {
		return NewDomainError("a negative amount %f can not be diposited", amount)
	}
	e := MoneyDipositedEvent{
		AccountId: a.accountId,
		EventId:   gocql.TimeUUID(),
		Amount:    amount,
	}
	if err := a.repo.Write(e); err != nil {
		return err
	}
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount < 0 {
		return NewDomainError("a negative amount %f can not be withdrawn", amount)
	}
	if a.balance-amount < a.limit {
		return NewDomainError("the withdrawn amount %f would exceed the limit", amount)
	}
	e := MoneyWithdrawnEvent{
		AccountId: a.accountId,
		EventId:   gocql.TimeUUID(),
		Amount:    amount,
	}
	if err := a.repo.Write(e); err != nil {
		return err
	}
	a.balance -= amount
	return nil
}

func (a *Account) Delete() error {
	e := AccountDeletedEvent{
		AccountId: a.accountId,
		EventId:   gocql.TimeUUID(),
	}
	if err := a.repo.Write(e); err != nil {
		return err
	}
	a.deleted = true
	return nil
}
