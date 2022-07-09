package domain

type EventStore interface {
	GetAllEvents() ([]*Event, error)
	InsertEvent(event *Event) error
}
