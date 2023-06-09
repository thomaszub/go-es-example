package domain

import (
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
		err := event.Apply(&acc)
		if err != nil {
			return acc, err
		}
	}
	if !acc.created || acc.deleted {
		return Account{}, NewAccountNotFoundError("account %s does not exist or is deleted", accountId)
	}
	return acc, nil
}

func (s *AccountService) GetAllAccountIds() ([]gocql.UUID, error) {
	activeIds := []gocql.UUID{}
	loadedIds, err := s.repo.ReadAllAccountIds()
	if err != nil {
		return activeIds, err
	}
	for _, id := range loadedIds {
		acc, err := s.GetAccount(id)
		if err != nil {
			return activeIds, err
		}
		if acc.Exists() {
			activeIds = append(activeIds, id)
		}
	}
	return activeIds, nil
}
