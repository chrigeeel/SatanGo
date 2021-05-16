package colors

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gookit/color"
)

func Prefix() string {
	dt := time.Now()
	dtf := dt.Format("15:04:05.000")
	red := color.FgLightRed.Render
	white := color.FgLightWhite.Render
	return white("["+dtf+"]") + red(" | ")
}
func TaskPrefix(id int) string {
	ids := strconv.Itoa(id)
	dt := time.Now()
	dtf := dt.Format("15:04:05.000")
	red := color.FgLightRed.Render
	white := color.FgLightWhite.Render
	return white("["+dtf+"]") + red(" | ") + white("[Task " + ids + "]") + red(" | ")

}
func Red(data string) string {
	red := color.FgLightRed.Render
	return fmt.Sprintf("%s", red(data))
}
func Green(data string) string {
	green := color.FgLightGreen.Render
	return fmt.Sprintf("%s", green(data))
}
func Yellow(data string) string {
	yellow := color.FgLightYellow.Render
	return fmt.Sprintf("%s", yellow(data))
}
func White(data string) string {
	white := color.FgLightWhite.Render
	return fmt.Sprintf("%s", white(data))
}
