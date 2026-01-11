// Package update provides types and interfaces for managing application updates.
package update

// List represents a collection of available updates.
type List struct {
	// Updates is the slice of available update items.
	Updates []Item

	// HasBlockingUpdate indicates if any update requires immediate action.
	HasBlockingUpdate bool
}

// Item represents a single update item.
type Item struct {
	// Name is the update identifier or package name.
	Name string

	// Version is the new version string.
	Version string

	// CurrentVersion is the currently installed version, if any.
	CurrentVersion string

	// IsBlocking indicates if this update blocks application usage.
	IsBlocking bool

	// Size is the download size in bytes, if known.
	Size int64

	// Description provides details about the update.
	Description string
}

// Interface defines the contract for update sources.
// Implementations provide version checking and update information
// for different components (launcher, game, JRE, etc.).
type Interface interface {
	// Name returns the identifier for this update source.
	Name() string

	// CheckForUpdate checks if an update is available.
	// Returns the new version info if available, nil if up-to-date.
	CheckForUpdate(channel string) (*Item, error)

	// Populate fills in additional details for an update item.
	// This may involve fetching download URLs, checksums, etc.
	Populate(item *Item, channel string) error
}

// IsEmpty returns true if the list contains no updates.
func (l *List) IsEmpty() bool {
	return len(l.Updates) == 0
}

// Add appends an update item to the list.
// If the item is blocking, HasBlockingUpdate is set to true.
func (l *List) Add(item Item) {
	l.Updates = append(l.Updates, item)
	if item.IsBlocking {
		l.HasBlockingUpdate = true
	}
}

// Clear removes all updates from the list.
func (l *List) Clear() {
	l.Updates = nil
	l.HasBlockingUpdate = false
}

// Event represents an update event emitted during the update process.
type Event struct {
	// Name is the event identifier.
	Name string `json:"name"`

	// Package is the package being updated.
	Package string `json:"package,omitempty"`

	// Version is the version being installed.
	Version string `json:"version,omitempty"`

	// Progress is the update progress (0-100).
	Progress float64 `json:"progress,omitempty"`

	// Error contains error details if the event represents a failure.
	Error string `json:"error,omitempty"`
}

// Notification represents a status update notification.
type Notification struct {
	// Package is the package being updated.
	Package string `json:"package,omitempty"`

	// Status is a human-readable status message.
	Status string `json:"status,omitempty"`

	// Progress is the overall progress (0-100).
	Progress float64 `json:"progress,omitempty"`

	// BytesDownloaded is the number of bytes downloaded so far.
	BytesDownloaded int64 `json:"bytes_downloaded,omitempty"`

	// BytesTotal is the total number of bytes to download.
	BytesTotal int64 `json:"bytes_total,omitempty"`

	// Speed is the current download speed in bytes per second.
	Speed int64 `json:"speed,omitempty"`
}

// Listener is an interface for receiving update events and notifications.
type Listener interface {
	// Event is called when an update event occurs.
	Event(event Event)

	// Notify is called to report update status.
	Notify(notification Notification)
}

// Schedule represents a scheduled update time window.
type Schedule struct {
	// StartTime is when the update window opens.
	StartTime string `json:"start_time,omitempty"`

	// EndTime is when the update window closes.
	EndTime string `json:"end_time,omitempty"`
}

// Package is an interface that update packages must implement.
type Package interface {
	// Name returns the package identifier.
	Name() string
}

// JREPackage represents the Java Runtime Environment update package.
type JREPackage struct{}

// Name returns "jre".
func (p *JREPackage) Name() string { return "jre" }

// GamePackage represents the game update package.
type GamePackage struct{}

// Name returns "game".
func (p *GamePackage) Name() string { return "game" }
