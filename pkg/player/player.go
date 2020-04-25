// Package player will be used control a player and get/manipulate playback.
package player

// EventCallback will be the type of function signature whenever an event occurs
type EventCallback func(data interface{}) error

// Player will be the interface which has to be implemented by any player to control the player
type Player interface {
	Listen() error
	PauseCallback(pauseCB EventCallback) error
	ExitCallback(exitCB EventCallback) error
	Close()
}
