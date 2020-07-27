package spinner

import (
	"fmt"
	"io"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/kaduartur/go-cli-spinner/pkg/template"
)

const (
	DefaultFrameRate = time.Millisecond * 150
)

type Spinner struct {
	sync.Mutex
	Title     string
	Template  template.Template
	FrameRate time.Duration
	runChan   chan struct{}
	stopOnce  sync.Once
	Output    io.Writer
	NoTty     bool
	timer     *time.Timer
}

// NewSpinner creates a Spinner with default template and FrameRate
func New(title string) *Spinner {
	sp := &Spinner{
		Title:     title,
		Template:  template.Default,
		FrameRate: DefaultFrameRate,
		runChan:   make(chan struct{}),
		timer:     time.NewTimer(DefaultFrameRate),
	}

	var stdout interface{} = syscall.Stdout
	stdoutFD, ok := stdout.(int)
	isTerminal := terminal.IsTerminal(stdoutFD)
	if !(ok && isTerminal) {
		sp.NoTty = true
	}

	return sp
}

// StartNew creates a Spinner with default template and FrameRate.
// after create Spinner, starts the execution.
func StartNew(title string) *Spinner {
	return New(title).Start()
}

// Start starts the spinner execution
func (s *Spinner) Start() *Spinner {
	go s.writer()
	return s
}

// SetSpeed sets the frame rate execution
func (s *Spinner) SetSpeed(rate time.Duration) *Spinner {
	s.Lock()
	s.FrameRate = rate
	s.Unlock()
	return s
}

// SetTemplate sets the spinner template
func (s *Spinner) SetTemplate(template template.Template) *Spinner {
	s.Lock()
	s.Template = template
	s.Unlock()
	return s
}

// Stop stops the spinner execution
func (s *Spinner) Stop() {
	s.stopOnce.Do(func() {
		close(s.runChan)
		if !s.timer.Stop() { // stop the timer when the spinner was stopped, preventing the next line from being cleared
			<-s.timer.C
		}
		s.clearLine()
	})
}

// Success gives success feedback for users and stops the spinner execution
func (s *Spinner) Success(msg string) {
	s.Stop()
	fmt.Println(msg)
}

// Error gives error feedback for users and stops the spinner execution
func (s *Spinner) Error(err error) {
	s.Stop()
	fmt.Println(err)
}

// animate runs the template animation
func (s *Spinner) animate() {
	for _, c := range s.Template {
		out := fmt.Sprintf("%s %s", c, s.Title)
		switch {
		case s.Output != nil:
			_, _ = fmt.Fprint(s.Output, out)
		case !s.NoTty:
			fmt.Print(out)
		}

		s.timer.Reset(s.FrameRate)
		<-s.timer.C // Wait for timer
		s.clearLine()
	}
}

// writer writes out spinner animation until runChan is closed
func (s *Spinner) writer() {
	for {
		select {
		case <-s.runChan:
			return
		default:
			s.animate()
		}
	}
}

func (s *Spinner) clearLine() {
	if !s.NoTty {
		fmt.Printf("\033[2K")
		fmt.Println()
		fmt.Printf("\033[1A")
	}
}
