package utility

import "github.com/chrigeeel/satango/loader"

func RemoveIndex(profiles []loader.ProfileStruct, s int) []loader.ProfileStruct {
	return append(profiles[:s], profiles[s+1:]...)
}