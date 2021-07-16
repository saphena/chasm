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
	DBH.Exec(chasmSQL)

	if LangCodeOverride == "" {
		return
	}

	sql := "UPDATE config SET Langcode=?"
	stmt, _ := DBH.Prepare(sql)
	defer stmt.Close()
	stmt.Exec(LangCodeOverride)

}
