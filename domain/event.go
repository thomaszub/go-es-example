package domain

import "github.com/gocql/gocql"

type AccountEvent interface {
	GetAccountId() gocql.UUID
	GetEventId() gocql.UUID
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
