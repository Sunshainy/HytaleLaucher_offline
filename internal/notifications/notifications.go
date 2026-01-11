// Package notifications provides system notification functionality
// for the Hytale launcher.
package notifications

import (
	"log/slog"
)

// Notification represents a system notification to be displayed to the user.
type Notification struct {
	// Title is the notification title.
	Title string `json:"title"`

	// Message is the notification message body.
	Message string `json:"message"`

	// Type indicates the notification type (info, warning, error, success).
	Type NotificationType `json:"type"`
}

// NotificationType represents the type/severity of a notification.
type NotificationType string

const (
	// TypeInfo is for informational notifications.
	TypeInfo NotificationType = "info"

	// TypeWarning is for warning notifications.
	TypeWarning NotificationType = "warning"

	// TypeError is for error notifications.
	TypeError NotificationType = "error"

	// TypeSuccess is for success notifications.
	TypeSuccess NotificationType = "success"
)

// Notifier is the interface for sending system notifications.
type Notifier interface {
	// Send displays a notification to the user.
	Send(n Notification) error
}

// logNotifier is a simple notifier that logs notifications.
type logNotifier struct{}

// defaultNotifier is the default notification handler.
var defaultNotifier Notifier = &logNotifier{}

// SetNotifier sets the default notifier implementation.
func SetNotifier(n Notifier) {
	defaultNotifier = n
}

// Send sends a notification using the default notifier.
func Send(n Notification) error {
	return defaultNotifier.Send(n)
}

// SendInfo sends an informational notification.
func SendInfo(title, message string) error {
	return Send(Notification{
		Title:   title,
		Message: message,
		Type:    TypeInfo,
	})
}

// SendWarning sends a warning notification.
func SendWarning(title, message string) error {
	return Send(Notification{
		Title:   title,
		Message: message,
		Type:    TypeWarning,
	})
}

// SendError sends an error notification.
func SendError(title, message string) error {
	return Send(Notification{
		Title:   title,
		Message: message,
		Type:    TypeError,
	})
}

// SendSuccess sends a success notification.
func SendSuccess(title, message string) error {
	return Send(Notification{
		Title:   title,
		Message: message,
		Type:    TypeSuccess,
	})
}

// Send implements Notifier by logging the notification.
func (l *logNotifier) Send(n Notification) error {
	slog.Info("notification",
		"title", n.Title,
		"message", n.Message,
		"type", n.Type,
	)
	return nil
}
