package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/thomaszub/go-es-example/domain"
)

type accountEventType string

const (
	accountCreatedEventType accountEventType = "created"
	accountDeletedEventType accountEventType = "deleted"
	moneyDipositedEventType accountEventType = "moneyDeposited"
	moneyWithdrawnEventType accountEventType = "moneyWithdrawn"
	limitSetEventType       accountEventType = "limitSet"
)

var accountEventTable = table.New(table.Metadata{
	Name:    "account_event",
	Columns: []string{"account_id", "event_id", "payload"},
	PartKey: []string{"account_id"},
	SortKey: []string{"event_id"},
})

type CqlAccountEventRepository struct {
	session gocqlx.Session
}

type PersistableAccountEvent struct {
	AccountId gocql.UUID
	EventId   gocql.UUID
	Payload   []byte
}

func InitRepository(session *gocql.Session) CqlAccountEventRepository {
	return CqlAccountEventRepository{
		session: gocqlx.NewSession(session),
	}
}

func (r *CqlAccountEventRepository) Write(event domain.AccountEvent) error {
	switch e := event.(type) {
	case domain.AccountCreatedEvent:
		return r.writeAccountCreatedEvent(e)
	case domain.AccountDeletedEvent:
		return r.writeAccountDeletedEvent(e)
	case domain.MoneyDipositedEvent:
		return r.writeMoneyDipositedEvent(e)
	case domain.MoneyWithdrawnEvent:
		return r.writeMoneyWithdrawnEvent(e)
	case domain.LimitSetEvent:
		return r.writeLimitSetEvent(e)
	default:
		return fmt.Errorf("%+v is not a valid account event", event)
	}
}

func (r *CqlAccountEventRepository) ReadAllEvents(accountId gocql.UUID) ([]domain.AccountEvent, error) {
	var loadedEvents []PersistableAccountEvent
	q := r.session.Query(accountEventTable.Select()).BindMap(qb.M{"account_id": accountId})
	if err := q.SelectRelease(&loadedEvents); err != nil {
		return []domain.AccountEvent{}, err
	}

	var deserEvents []domain.AccountEvent
	for _, event := range loadedEvents {
		payload := make(map[string]interface{})
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return deserEvents, err
		}
		typeI, ok := payload["eventType"]
		if !ok {
			return deserEvents, fmt.Errorf("type is not set on event %s", event.EventId.String())
		}
		eventTypeF, ok := typeI.(string)
		if !ok {
			return deserEvents, fmt.Errorf("%v is not a valid event type for event %s", typeI, event.EventId.String())
		}
		eventType := accountEventType(eventTypeF)
		switch eventType {
		case accountCreatedEventType:
			e := domain.AccountCreatedEvent{
				AccountId: event.AccountId,
				EventId:   event.EventId,
			}
			deserEvents = append(deserEvents, e)
		case accountDeletedEventType:
			e := domain.AccountDeletedEvent{
				AccountId: event.AccountId,
				EventId:   event.EventId,
			}
			deserEvents = append(deserEvents, e)
		case moneyDipositedEventType:
			e, err := deserializeMoneyDipositedEvent(event.AccountId, event.EventId, payload)
			if err != nil {
				return deserEvents, errors.Join(fmt.Errorf("event %s", event.EventId.String()), err)
			}
			deserEvents = append(deserEvents, e)
		case moneyWithdrawnEventType:
			e, err := deserializeMoneyWithdrawnEvent(event.AccountId, event.EventId, payload)
			if err != nil {
				return deserEvents, errors.Join(fmt.Errorf("event %s", event.EventId.String()), err)
			}
			deserEvents = append(deserEvents, e)
		case limitSetEventType:
			e, err := deserializeLimitSetEvent(event.AccountId, event.EventId, payload)
			if err != nil {
				return deserEvents, errors.Join(fmt.Errorf("event %s", event.EventId.String()), err)
			}
			deserEvents = append(deserEvents, e)

		default:
			return deserEvents, fmt.Errorf("%s is not a known event type", eventType)
		}
	}
	return deserEvents, nil
}

func (r *CqlAccountEventRepository) ReadAllAccountIds() ([]gocql.UUID, error) {
	var ids []gocql.UUID
	q := r.session.Query(qb.Select(accountEventTable.Metadata().Name).Columns("account_id").Distinct("account_id").ToCql())
	if err := q.SelectRelease(&ids); err != nil {
		return ids, err
	}
	return ids, nil
}

func (r *CqlAccountEventRepository) writeAccountCreatedEvent(event domain.AccountCreatedEvent) error {
	payload := fmt.Sprintf(`{"eventType":"%s"}`, accountCreatedEventType)
	return r.write(event.AccountId, event.EventId, payload)
}

func (r *CqlAccountEventRepository) writeAccountDeletedEvent(event domain.AccountDeletedEvent) error {
	payload := fmt.Sprintf(`{"eventType":"%s"}`, accountDeletedEventType)
	return r.write(event.AccountId, event.EventId, payload)
}

func (r *CqlAccountEventRepository) writeMoneyDipositedEvent(event domain.MoneyDipositedEvent) error {
	payload := fmt.Sprintf(`{"eventType":"%s","amount":%f}`, moneyDipositedEventType, event.Amount)
	return r.write(event.AccountId, event.EventId, payload)
}

func (r *CqlAccountEventRepository) writeMoneyWithdrawnEvent(event domain.MoneyWithdrawnEvent) error {
	payload := fmt.Sprintf(`{"eventType":"%s","amount":%f}`, moneyWithdrawnEventType, event.Amount)
	return r.write(event.AccountId, event.EventId, payload)
}

func (r *CqlAccountEventRepository) writeLimitSetEvent(event domain.LimitSetEvent) error {
	payload := fmt.Sprintf(`{"eventType":"%s","limit":%f}`, limitSetEventType, event.Limit)
	return r.write(event.AccountId, event.EventId, payload)
}

func (r *CqlAccountEventRepository) write(accountId, eventId gocql.UUID, payload string) error {
	pe := PersistableAccountEvent{
		AccountId: accountId,
		EventId:   eventId,
		Payload:   []byte(payload),
	}
	return r.session.Query(accountEventTable.Insert()).BindStruct(pe).ExecRelease()
}

func deserializeMoneyDipositedEvent(accountId, eventId gocql.UUID, payload map[string]interface{}) (domain.MoneyDipositedEvent, error) {
	e := domain.MoneyDipositedEvent{
		AccountId: accountId,
		EventId:   eventId,
	}
	amount, err := getTypedValue[float64](payload, "amount")
	if err != nil {
		return e, err
	}
	e.Amount = amount
	return e, nil
}

func deserializeMoneyWithdrawnEvent(accountId, eventId gocql.UUID, payload map[string]interface{}) (domain.MoneyWithdrawnEvent, error) {
	e := domain.MoneyWithdrawnEvent{
		AccountId: accountId,
		EventId:   eventId,
	}
	amount, err := getTypedValue[float64](payload, "amount")
	if err != nil {
		return e, err
	}
	e.Amount = amount
	return e, nil
}

func deserializeLimitSetEvent(accountId, eventId gocql.UUID, payload map[string]interface{}) (domain.LimitSetEvent, error) {
	e := domain.LimitSetEvent{
		AccountId: accountId,
		EventId:   eventId,
	}
	limit, err := getTypedValue[float64](payload, "limit")
	if err != nil {
		return e, err
	}
	e.Limit = limit
	return e, nil
}

func getTypedValue[T any](payload map[string]interface{}, key string) (T, error) {
	var value T
	valueS, ok := payload[key]
	if !ok {
		return value, fmt.Errorf("%s is not set", key)
	}
	value, ok = valueS.(T)
	if !ok {
		return value, fmt.Errorf("%v is not of type %T", value, value)
	}
	return value, nil
}
