package discordjoiner

import (
	"fmt"
	"strconv"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/utility"
)


func Input(userData loader.UserDataStruct, profiles []loader.ProfileStruct, mode string) {
	profiles = utility.AskForProfiles(profiles)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Prefix() + colors.Red("Running ") + colors.White(strconv.Itoa(len(profiles)) + colors.Red(" tasks.")))
}