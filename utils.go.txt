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

func createDatabase(LangCodeOverride string) {

	fmt.Printf("Establishing initial database structure\n")
	defer DBH.Close()

	_, err := DBH.Exec(chasmSQL)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	if LangCodeOverride == "" {
		return
	}

	sql := "UPDATE config SET Langcode=?"
	_, err = DBH.Exec(sql, LangCodeOverride)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
