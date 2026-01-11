package app

import (
	"hytale-launcher/internal/update"
)

// appListen implements the update.Listener interface.
// It forwards update events and notifications to the frontend via the App's Emit method.
type appListen struct {
	// emit is a function that sends events to the frontend.
	emit func(name string, args ...any)
}

// Event forwards an update event to the frontend.
// It wraps the event in a slice and emits it with the event name.
func (l *appListen) Event(event update.Event) {
	l.emit(event.Name, event)
}

// Notify forwards an update notification to the frontend.
// Notifications are typically used for status updates during downloads/updates.
func (l *appListen) Notify(notification update.Notification) {
	l.emit("update:status", notification)
}

// newAppListen creates a new appListen instance with the given emit function.
func newAppListen(emit func(name string, args ...any)) *appListen {
	return &appListen{emit: emit}
}
