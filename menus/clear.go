package menus

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/chrigeeel/satango/colors"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
	var toPrint string
	toPrint += "       .-._                                                   _,-,\n"
	toPrint += "        `._`-._                                           _,-'_,'\n"
	toPrint += "            `._ `-._                                   _,-' _,'\n"
	toPrint += "                `._  `-._        __.-----.__        _,-'  _,'\n"
	toPrint += "                `._   `#===\"\"\"           \"\"\"===#'   _,'\n"
	toPrint += "                    `._/)  ._               _.  (\\_,'\n"
	toPrint += "                    )*'     **.__     __.**     '*(\n"
	toPrint += "                    #  .==..__  \"\"   \"\"  __..==,  #\n"
	toPrint += "                    #   `\"._(_).       .(_)_.\"'   #\n"
	toPrint += "_____       __      #       ____        __       _#______    ____\n"
	toPrint += "/ ___/____ _/ /_____ #____  / __ )____  / /_     / ____/ /   /  _/\n"
	toPrint += "\\__ \\/ __ `/ __/ __ `/ __ \\/ __  / __ \\/ __/    / /   / /    / /\n"
	toPrint += "___/ / /_/ / /_/ /_/ / / / / /_/ / /_/ / /_     / /___/ /____/ /\n"
	toPrint += "/____/\\__,_/\\__/\\__,_/_/ /_/_____/\\____/\\__/     \\____/_____/___/"
	fmt.Println(colors.Red(toPrint))
	fmt.Println(colors.Red("------------------------------------------------------------------"))
}
