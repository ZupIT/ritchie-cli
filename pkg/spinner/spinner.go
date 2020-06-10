package spinner

import (
	"time"

	"github.com/janeczku/go-spinner"
)

var (
	speed     = 100 * time.Millisecond
	sleepTime = 3 * time.Second
	template  = []string{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"}
)

// Manager is a wrapper of the go-spinner library,
// that allows us to standardize spinners on Ritchie-CLI
type Manager struct {
	speed     time.Duration
	sleepTime time.Duration
	Spinner   *spinner.Spinner
}

// New creates a new "spinner" and set the message
// to the "spinner" with a standard template
func New(startMgs string) *Manager {
	s := spinner.NewSpinner(startMgs)
	s.SetCharset(template)
	return &Manager{
		Spinner:   s,
		speed:     speed,
		sleepTime: sleepTime,
	}
}

// Start sets spinner speed, start spinner
// and wait by the time set in 'sleepTime' variable
func (m *Manager) Start() {
	m.Spinner.SetSpeed(m.speed)
	m.Spinner.Start()
	time.Sleep(m.sleepTime)
}

// Stop stops spinner and print the stop message
func (m *Manager) Stop() {
	m.Spinner.Stop()
}

func (m *Manager) SetSleepTime(d time.Duration) {
	m.sleepTime = d
}

func (m *Manager) SetSpeed(d time.Duration) {
	m.speed = d
}
