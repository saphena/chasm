package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func getYN(prompt string) bool {

	promptx := promptui.Select{
		Label: prompt,
		Items: []string{"Yes", "No"},
	}

	_, result, _ := promptx.Run()

	fmt.Printf("You chose %v\n", result)
	return result == "Yes"

}

func createDatabase() {

	fmt.Printf("Establishing initial database structure\n")
	DBH.Exec(chasmSQL)
}
