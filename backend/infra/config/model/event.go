package model

type EventType int

const (
	Update EventType = iota
	Delete
	Create
	Unknown
)

type Event struct {
	ConfigItem
	EventType
}

type EventHandler func([]*Event)

func (s *Event) String() string {
	out := "unknown"
	switch s.EventType {
	case Create:
		out = "create"
	case Update:
		out = "update"
	case Delete:
		out = "delete"
	}
	return out
}
