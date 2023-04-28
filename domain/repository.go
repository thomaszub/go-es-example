package domain

import (
	"github.com/gocql/gocql"
)

type AccountEventRepository interface {
	Write(event AccountEvent) error
	ReadAllEvents(accountId gocql.UUID) ([]AccountEvent, error)
	ReadAllAccountIds() ([]gocql.UUID, error)
}
