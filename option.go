package pomato

import "time"

type Option func(p *Pomato)

func WithPomodoroTime(d time.Duration) Option {
	return func(p *Pomato) {
		p.pomodoroTime = d
	}
}

func WithBreakTime(d time.Duration) Option {
	return func(p *Pomato) {
		p.breakTime = d
	}
}

func WithLongBreakTime(d time.Duration) Option {
	return func(p *Pomato) {
		p.longBreakTime = d
	}
}

func WithLongBreakEach(rounds int) Option {
	return func(p *Pomato) {
		p.longBreakEach = rounds
	}
}

func WithAutoStartNext(auto bool) Option {
	return func(p *Pomato) {
		p.autoStartNext = auto
	}
}

func WithShowNotification(show bool) Option {
	return func(p *Pomato) {
		p.showNotification = show
	}
}
