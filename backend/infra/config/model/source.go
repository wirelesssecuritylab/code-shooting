package model

type ConfigSource interface {
	Store
	Monitor
	GetSourceName() string
}
