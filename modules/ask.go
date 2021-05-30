package modules

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
)

func askForSilent() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(colors.Prefix() + colors.White("> "))
	scanner.Scan()
	return scanner.Text()
}

func askForProfiles(profiles []loader.ProfileStruct) []loader.ProfileStruct {
	fmt.Println(colors.Prefix() + colors.Red("What profiles do you want to run? You have the following profiles:"))
	for i := range profiles {
		fmt.Println(colors.Prefix() + colors.White(""+strconv.Itoa(i+1)+") \"") + colors.Red(profiles[i].Name) + colors.White("\""))
	}
	fmt.Println(colors.Prefix() + colors.Red("Sample input: ") + colors.White("\"") + colors.Red("1, 2, 5") + colors.White("\" ") + colors.Red("or ") + colors.White("\"") + colors.Red("all") + colors.White("\""))
	var newProfiles []loader.ProfileStruct
	var indexes []int
	for validAns := false; validAns == false; {
		ans := askForSilent()
		validAns = true
		if strings.ToLower(ans) == "all" {
			for i := range profiles {
				indexes = append(indexes, i)
			}
			break
		}
		ans = strings.ReplaceAll(ans, " ", "")
		plist1 := strings.Split(ans, ",")
		for i := range plist1 {
			if govalidator.IsInt(plist1[i]) == false {
				fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
				validAns = false
			} else {
				pint, err := strconv.Atoi(plist1[i])
				if err != nil {
					fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
					validAns = false
				}
				if pint <= len(profiles) {
					pint--
					indexes = append(indexes, pint)
				} else {
					fmt.Println(colors.Prefix() + colors.Red("You don't have the profile " + strconv.Itoa(pint)))
					validAns = false
				}
			}
		}
	}

	var toPrint string
	toPrint += colors.Prefix() + colors.Red("Running the profiles: \n") + colors.Prefix() 
	for i := range indexes {
		newProfiles = append(newProfiles, profiles[indexes[i]])
		toPrint = toPrint + colors.White("\"") + colors.Red(profiles[indexes[i]].Name) + colors.White("\", ")
	}
	fmt.Println(toPrint)
	return newProfiles
}