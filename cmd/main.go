package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alfred-zhong/pomato"
	"github.com/spf13/viper"
)

var version string

func main() {
	// register config
	registerViper()
	// read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// no config file found.
		} else {
			fmt.Printf("read in config file: %v\n", err)
			os.Exit(1)
		}
	}

	// create pomato and run
	p := pomato.NewPomato(buildOptions()...)
	if err := p.Run(); err != nil {
		fmt.Printf("pomato run fail: %v\n", err)
		os.Exit(2)
	}
}

var (
	configName = "pomato"
	configType = "yaml"
	configPath = []string{
		"/etc", "$HOME", ".",
	}
)

func registerViper() {
	viper.SetConfigName(configName)
	// viper.SetConfigType(configType)
	for _, path := range configPath {
		viper.AddConfigPath(path)
	}
}

const (
	keyPomodoroTime     = "pomodoro_time"
	keyBreakTime        = "break_time"
	keyLongBreakTime    = "long_break_time"
	keyLongBreakEach    = "long_break_each"
	keyAutostartNext    = "autostart_next"
	keyShowNotification = "show_notification"
)

func buildOptions() []pomato.Option {
	opts := make([]pomato.Option, 0, 6)
	if i := viper.GetInt(keyPomodoroTime); i > 0 {
		opts = append(opts, pomato.WithPomodoroTime(time.Duration(i)*time.Minute))
	}
	if i := viper.GetInt(keyBreakTime); i > 0 {
		opts = append(opts, pomato.WithBreakTime(time.Duration(i)*time.Minute))
	}
	if i := viper.GetInt(keyLongBreakTime); i > 0 {
		opts = append(opts, pomato.WithLongBreakTime(time.Duration(i)*time.Minute))
	}
	if r := viper.GetInt(keyLongBreakEach); r > 0 {
		opts = append(opts, pomato.WithLongBreakEach(r))
	}
	if viper.InConfig(keyAutostartNext) {
		opts = append(opts, pomato.WithAutoStartNext(viper.GetBool(keyAutostartNext)))
	}
	if viper.InConfig(keyShowNotification) {
		opts = append(opts, pomato.WithShowNotification(viper.GetBool(keyShowNotification)))
	}
	return opts
}
