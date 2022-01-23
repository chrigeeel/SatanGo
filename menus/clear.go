package menus

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/chrigeeel/satango/colors"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] 
	if ok {                  
		value()
	} else {
		for i := 0; i < 50; i++ {
			fmt.Println("")
		}
	}
	var toPrint string
	toPrint += "          .-._                                                   _,-,\n"
	toPrint += "           `._`-._                                           _,-'_,'\n"
	toPrint += "               `._ `-._                                   _,-' _,'\n"
	toPrint += "                   `._  `-._        __.-----.__        _,-'  _,'\n"
	toPrint += "                   `._   `#===\"\"\"           \"\"\"===#'   _,'\n"
	toPrint += "                       `._/)  ._               _.  (\\_,'\n"
	toPrint += "                       )*'     **.__     __.**     '*(\n"
	toPrint += "                       #  .==..__  \"\"   \"\"  __..==,  #\n"
	toPrint += "                       #   `\"._(_).       .(_)_.\"'   #\n"
	toPrint += "   _____       __      #       ____        __       _#______    ____\n"
	toPrint += "  / ___/____ _/ /_____ #____  / __ )____  / /_     / ____/ /   /  _/\n"
	toPrint += "  \\__ \\/ __ `/ __/ __ `/ __ \\/ __  / __ \\/ __/    / /   / /    / /\n"
	toPrint += " ___/ / /_/ / /_/ /_/ / / / / /_/ / /_/ / /_     / /___/ /____/ /\n"
	toPrint += "/____/\\__,_/\\__/\\__,_/_/ /_/_____/\\____/\\__/     \\____/_____/___/"
	fmt.Println(colors.Red(toPrint))
	fmt.Println(colors.Red("------------------------------------------------------------------"))
}
