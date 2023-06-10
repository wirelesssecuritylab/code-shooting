package model

type Monitor interface {
	Start(stop <-chan struct{}) error
	RegisterEventHandler(key string, handler EventHandler)
	ProcessConfigEvent(events []*Event)
}
