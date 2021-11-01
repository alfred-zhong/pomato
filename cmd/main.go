package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/alfred-zhong/pomato"
	"github.com/spf13/pflag"
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

	// bind with flag
	flag.Int(keyPomodoroTime, pomato.DefaultPomodoroTime, "pomodoro time. (Unit: minute)")
	flag.Int(keyBreakTime, pomato.DefaultBreakTime, "break time. (Unit: minute)")
	flag.Int(keyLongBreakTime, pomato.DefaultLongBreakTime, "long break time. (Unit: minute)")
	flag.Int(keyLongBreakEach, pomato.DefaultLongBreakEach, "long break each rounds.")
	flag.String(keyTimeUnit, "m", "time unit. (\"m\" or \"s\")")
	flag.Bool(keyAutostartNext, pomato.DefaultAutoStartNext, "auto start next pomodoro.")
	flag.Bool(keyShowNotification, pomato.DefaultShowNotification, "show notification.")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

const (
	keyPomodoroTime     = "pomodoro-time"
	keyBreakTime        = "break-time"
	keyLongBreakTime    = "long-break-time"
	keyLongBreakEach    = "long-break-each"
	keyTimeUnit         = "time-unit"
	keyAutostartNext    = "autostart-next"
	keyShowNotification = "show-notification"
)

func buildOptions() []pomato.Option {
	opts := make([]pomato.Option, 0, 6)
	timeUnit := time.Minute
	if s := viper.GetString(keyTimeUnit); s == "s" {
		timeUnit = time.Second
	}
	if i := viper.GetInt(keyPomodoroTime); i > 0 {
		opts = append(opts, pomato.WithPomodoroTime(time.Duration(i)*timeUnit))
	}
	if i := viper.GetInt(keyBreakTime); i > 0 {
		opts = append(opts, pomato.WithBreakTime(time.Duration(i)*timeUnit))
	}
	if i := viper.GetInt(keyLongBreakTime); i > 0 {
		opts = append(opts, pomato.WithLongBreakTime(time.Duration(i)*timeUnit))
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
