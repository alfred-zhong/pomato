package pomato

import (
	"fmt"
	"os/exec"
)

func notifyStart(title, text string) {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s" sound name "Blow.aiff"`, text, title))
	fmt.Println(cmd)
	if err := cmd.Run(); err != nil {
		fmt.Printf("notify fail: %v", err)
	}
}

func notifyEnd(title, text string) {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s" sound name "Bottle.aiff"`, text, title))
	if err := cmd.Run(); err != nil {
		fmt.Printf("notify fail: %v", err)
	}
}
