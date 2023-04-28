package domain

import (
	"fmt"

	"github.com/gocql/gocql"
)

type AccountService struct {
	repo AccountEventRepository
}

func NewAccountService(repo AccountEventRepository) AccountService {
	return AccountService{
		repo: repo,
	}
}

func (s *AccountService) GetAllAccountIds() ([]gocql.UUID, error) {
	var activeIds []gocql.UUID
	loadedIds, err := s.repo.ReadAllAccountIds()
	if err != nil {
		return activeIds, err
	}
	for _, id := range loadedIds {
		exists, err := s.accountExists(id)
		if err != nil {
			return activeIds, err
		}
		if exists {
			activeIds = append(activeIds, id)
		}
	}
	return activeIds, nil
}

func (s *AccountService) CreateNewAccount() (Account, error) {
	e := AccountCreatedEvent{
		AccountId: gocql.MustRandomUUID(),
		EventId:   gocql.TimeUUID(),
	}
	if err := s.repo.Write(e); err != nil {
		return Account{}, err
	}
	return s.GetAccount(e.AccountId)
}

func (s *AccountService) GetAccount(accountId gocql.UUID) (Account, error) {
	acc := Account{
		repo:      s.repo,
		accountId: accountId,
		created:   false,
		deleted:   false,
	}
	events, err := s.repo.ReadAllEvents(accountId)
	if err != nil {
		return acc, err
	}
	for _, event := range events {
		switch e := event.(type) {
		case AccountCreatedEvent:
			acc.created = true
		case AccountDeletedEvent:
			acc.deleted = true
		case MoneyDipositedEvent:
			acc.balance += e.Amount
		case MoneyWithdrawnEvent:
			acc.balance -= e.Amount
		case LimitSetEvent:
			acc.limit = e.Limit
		default:
			return acc, fmt.Errorf("%+v is not a valid account event", event)
		}
	}
	if !acc.created || acc.deleted {
		return Account{}, NewAccountNotFoundError("account %s does not exist or is deleted", accountId)
	}
	return acc, nil
}

func (s *AccountService) accountExists(accountId gocql.UUID) (bool, error) {
	events, err := s.repo.ReadAllEvents(accountId)
	if err != nil {
		return false, err
	}
	exists := false
	for _, event := range events {
		switch event.(type) {
		case AccountCreatedEvent:
			exists = true
		case AccountDeletedEvent:
			exists = false
		}
	}
	return exists, nil
}
