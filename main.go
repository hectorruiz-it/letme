package main

import (
	"fmt"

	utils "github.com/lockedinspace/letme/pkg"
	letme "github.com/lockedinspace/letme/pkg/cmd"
)

func main() {
	if utils.CacheFileExists() {
		fmt.Println("letme: file" + utils.GetHomeDirectory() + "/.letme/.letme-cache" + " not supported anymore. Please remove it manually.")
	}
	utils.CommandExists("aws")
	letme.Execute()
}
