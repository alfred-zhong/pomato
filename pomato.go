package pomato

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// const (
// 	defaultPomodoroTime     = 6 * time.Second
// 	defaultBreakTime        = 3 * time.Second
// 	defaultLongBreakTime    = 5 * time.Second
// 	defaultLongBreakEach    = 2
// 	defaultAutoStartNext    = false
// 	defaultShowNotification = true
// )

const (
	defaultPomodoroTime     = 25 * time.Minute
	defaultBreakTime        = 5 * time.Minute
	defaultLongBreakTime    = 15 * time.Minute
	defaultLongBreakEach    = 4
	defaultAutoStartNext    = false
	defaultShowNotification = true
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
	for {
		for i := 0; i < p.longBreakEach; i++ {
			// Pomodoro count down
			countDown(
				fmt.Sprintf("[%d] %s", i+1, color.YellowString("Pomodoro time")),
				p.pomodoroTime,
			)

			// break count down
			if i == p.longBreakEach-1 {
				if p.showNotification {
					notifyStart("Pomato", "Relex and have a long break!")
				}
				countDown(
					fmt.Sprintf("[%d] %s", i+1, color.GreenString("Long break time")),
					p.longBreakTime,
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
					p.breakTime,
				)
				if p.showNotification {
					notifyEnd("Pomato", "Break finished!")
				}

				// check if need to wait until user press to continue
				if !p.autoStartNext {
					if _, err := waitForPress("Press enter to continue"); err != nil {
						return err
					}
				}
			}
		}

		// user must press to start a new big pomodoro
		if _, err := waitForPress("Press enter to start a new round"); err != nil {
			return err
		}
	}
}

func countDown(msg string, d time.Duration) error {
	endTime := time.Now().Add(d)

	timer := time.NewTimer(d)
	defer timer.Stop()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// print immediately
	var (
		left  = int(math.Round(time.Until(endTime).Seconds()))
		total = int(math.Round(d.Seconds()))
	)
	totalStr := fmt.Sprintf("%d:%02d", total/60, total%60)
	fmt.Printf("\r\033[K%s: %d:%02d/%s", msg, left/60, left%60, totalStr)
LOOP:
	for {
		select {
		case <-timer.C:
			fmt.Printf("\r\033[K%s: 0:00/%s", msg, totalStr)
			break LOOP
		case <-ticker.C:
			left = int(math.Round(time.Until(endTime).Seconds()))
			fmt.Printf("\r\033[K%s: %d:%02d/%s", msg, left/60, left%60, totalStr)
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

func NewPomato(options ...Option) *Pomato {
	p := &Pomato{
		pomodoroTime:     defaultPomodoroTime,
		breakTime:        defaultBreakTime,
		longBreakTime:    defaultLongBreakTime,
		longBreakEach:    defaultLongBreakEach,
		autoStartNext:    defaultAutoStartNext,
		showNotification: defaultShowNotification,
	}

	for _, opts := range options {
		opts(p)
	}
	return p
}
