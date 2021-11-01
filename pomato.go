package pomato

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	DefaultPomodoroTime     = 25
	DefaultBreakTime        = 5
	DefaultLongBreakTime    = 15
	DefaultLongBreakEach    = 4
	DefaultAutoStartNext    = false
	DefaultShowNotification = true
)

type Pomato struct {
	pomodoroTime     time.Duration
	breakTime        time.Duration
	longBreakTime    time.Duration
	longBreakEach    int
	autoStartNext    bool
	showNotification bool
}

func (p *Pomato) Run() error {
	return p.pomodoroMode()
}

func (p *Pomato) pomodoroMode() error {
	stopC := make(chan struct{})
	defer close(stopC)

	stdinC := listenStdin(stopC)

	for {
		for i := 0; i < p.longBreakEach; i++ {
			// Pomodoro count down
			countDown(
				fmt.Sprintf("[%d] %s", i+1, color.YellowString("Pomodoro time")),
				p.pomodoroTime, stdinC,
			)

			// break count down
			if i == p.longBreakEach-1 {
				if p.showNotification {
					notifyStart("Pomato", "Relex and have a long break!")
				}
				countDown(
					fmt.Sprintf("[%d] %s", i+1, color.GreenString("Long break time")),
					p.longBreakTime, stdinC,
				)
				if p.showNotification {
					notifyStart("Pomato", "Long break finished! Continue to work.")
				}
			} else {
				if p.showNotification {
					notifyStart("Pomato", "Have a break!")
				}
				countDown(
					fmt.Sprintf("[%d] %s", i+1, color.BlueString("Break time")),
					p.breakTime, stdinC,
				)
				if p.showNotification {
					notifyEnd("Pomato", "Break finished!")
				}

				// check if need to wait until user press to continue
				if !p.autoStartNext {
					waitForReturn("Press enter to continue", stdinC)
				}
			}
		}

		// user must press to start a new big pomodoro
		waitForReturn("Press enter to continue", stdinC)
	}
}

func countDown(msg string, d time.Duration, stdinC chan byte) error {
	endTime := time.Now().Add(d)

	stopC := make(chan struct{})
	defer close(stopC)

	spaceC := listenForSpace(stopC, stdinC)

	// timer := time.NewTimer(d)
	// defer timer.Stop()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// print immediately
	var (
		left  = int(math.Round(time.Until(endTime).Seconds()))
		total = int(math.Round(d.Seconds()))
	)
	totalStr := fmt.Sprintf("%d:%02d", total/60, total%60)
	fmt.Printf("\r\033[K%s: %d:%02d/%s", msg, left/60, left%60, totalStr)

	var (
		paused      bool
		leftSeconds float64
	)
LOOP:
	for {
		select {
		// case <-timer.C:
		// 	fmt.Printf("\r\033[K%s: 0:00/%s", msg, totalStr)
		// 	break LOOP
		case <-ticker.C:
			leftSeconds = time.Until(endTime).Seconds()
			left = int(math.Round(leftSeconds))
			fmt.Printf("\r\033[K%s: %d:%02d/%s", msg, left/60, left%60, totalStr)

			if left <= 0 {
				break LOOP
			}
		case <-spaceC:
			if paused {
				// timer.Reset(time.Duration(leftSeconds))
				// timer = time.NewTimer(time.Duration(leftSeconds))
				ticker.Reset(time.Second)

				endTime = time.Now().Add(time.Duration(leftSeconds) * time.Second)
				left = int(math.Round(leftSeconds))
				fmt.Printf("\r\033[K%s: %d:%02d/%s", msg, left/60, left%60, totalStr)
			} else {
				// timer.Stop()
				ticker.Stop()

				leftSeconds = time.Until(endTime).Seconds()
				left = int(math.Round(leftSeconds))
				fmt.Printf("\r\033[K%s: %d:%02d/%s [Paused]", msg, left/60, left%60, totalStr)
			}
			paused = !paused
		}
	}
	fmt.Println()

	return nil
}

func waitForPress(msg string) (input string, err error) {
	// discard all input before
	flushStdin()

	fmt.Printf("\r\033[K%s: ", msg)

	// read input
	r := bufio.NewReader(os.Stdin)
	b, err := r.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func waitForReturn(msg string, stdinC chan byte) {
	fmt.Printf("\r\033[K%s: ", msg)

	for b := range stdinC {
		if b == '\n' {
			return
		}
	}
}

func listenForSpace(stopC chan struct{}, stdinC chan byte) chan struct{} {
	c := make(chan struct{})

	go func() {
		defer close(c)

	LOOP:
		for {
			select {
			case <-stopC:
				break LOOP
			case b := <-stdinC:
				if b == ' ' {
					select {
					case c <- struct{}{}:
					default:
					}
				}
			}
		}
	}()

	return c
}

func listenStdin(stopC chan struct{}) chan byte {
	c := make(chan byte)

	go func() {
		defer close(c)

		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		buf := make([]byte, 1)
	LOOP:
		for {
			if _, err := os.Stdin.Read(buf); err != nil {
				fmt.Printf("read from stdin fail: %v\n", err)
				return
			}

			select {
			case <-stopC:
				break LOOP
			case c <- buf[0]:
			default:
			}
		}
	}()

	return c
}

func NewPomato(options ...Option) *Pomato {
	p := &Pomato{
		pomodoroTime:     DefaultPomodoroTime * time.Minute,
		breakTime:        DefaultBreakTime * time.Minute,
		longBreakTime:    DefaultLongBreakTime * time.Minute,
		longBreakEach:    DefaultLongBreakEach,
		autoStartNext:    DefaultAutoStartNext,
		showNotification: DefaultShowNotification,
	}

	for _, opts := range options {
		opts(p)
	}
	return p
}
