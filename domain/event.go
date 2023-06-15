package domain

import (
	"fmt"

	"github.com/gocql/gocql"
)

type AccountEvent interface {
	GetAccountId() gocql.UUID
	GetEventId() gocql.UUID
	Apply(account *Account) error
}

type AccountCreatedEvent struct {
	AccountId gocql.UUID
	EventId   gocql.UUID
}

func (e AccountCreatedEvent) GetAccountId() gocql.UUID {
	return e.AccountId
}

func (e AccountCreatedEvent) GetEventId() gocql.UUID {
	return e.EventId
}

func (e AccountCreatedEvent) Apply(account *Account) error {
	if e.AccountId != account.accountId {
		return eventAccountMismatched(e, account)
	}
	return nil
}

type AccountDeletedEvent struct {
	AccountId gocql.UUID
	EventId   gocql.UUID
}

func (e AccountDeletedEvent) GetAccountId() gocql.UUID {
	return e.AccountId
}

func (e AccountDeletedEvent) GetEventId() gocql.UUID {
	return e.EventId
}

func (e AccountDeletedEvent) Apply(account *Account) error {
	if e.AccountId != account.accountId {
		return eventAccountMismatched(e, account)
	}
	account.deleted = true
	return nil
}

type MoneyDipositedEvent struct {
	AccountId gocql.UUID
	EventId   gocql.UUID
	Amount    float64
}

func (e MoneyDipositedEvent) GetAccountId() gocql.UUID {
	return e.AccountId
}

func (e MoneyDipositedEvent) GetEventId() gocql.UUID {
	return e.EventId
}

func (e MoneyDipositedEvent) Apply(account *Account) error {
	if e.AccountId != account.accountId {
		return eventAccountMismatched(e, account)
	}
	account.balance += e.Amount
	return nil
}

type MoneyWithdrawnEvent struct {
	AccountId gocql.UUID
	EventId   gocql.UUID
	Amount    float64
}

func (e MoneyWithdrawnEvent) GetAccountId() gocql.UUID {
	return e.AccountId
}

func (e MoneyWithdrawnEvent) GetEventId() gocql.UUID {
	return e.EventId
}

func (e MoneyWithdrawnEvent) Apply(account *Account) error {
	if e.AccountId != account.accountId {
		return eventAccountMismatched(e, account)
	}
	account.balance -= e.Amount
	return nil
}

type LimitSetEvent struct {
	AccountId gocql.UUID
	EventId   gocql.UUID
	Limit     float64
}

func (e LimitSetEvent) GetAccountId() gocql.UUID {
	return e.AccountId
}

func (e LimitSetEvent) GetEventId() gocql.UUID {
	return e.EventId
}

func (e LimitSetEvent) Apply(account *Account) error {
	if e.AccountId != account.accountId {
		return eventAccountMismatched(e, account)
	}
	account.limit = e.Limit
	return nil
}

func eventAccountMismatched(event AccountEvent, account *Account) error {
	return fmt.Errorf("event %+v is not an event of account %s", event, account.accountId)
}
