package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

//go:embed chasmdb.sql
var chasmsql string

func ensureDirWritable(path string) error {
	// Ensure the path is absolute
	dir, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if directory exists
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// Create directory with 0755 permissions
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat directory: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", dir)
	}

	// Check if directory is writable by trying to create a temp file
	f, err := os.CreateTemp(dir, ".writable_check_*")
	if err != nil {
		return fmt.Errorf("directory is not writable: %w", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	return nil
}

func establishImageFolders() {

	fmt.Println("Checking/establishing image folders")
	err := ensureDirWritable(CS.ImgBonusFolder)
	if err != nil {
		fmt.Printf("Bonus image folder : %v\n", err)
	}
	err = ensureDirWritable(CS.ImgEbcFolder)
	if err != nil {
		fmt.Printf("EBC image folder : %v\n", err)
	}
}

func getIntegerFromDB(sqlx string, defval int) int {

	str := getStringFromDB(sqlx, strconv.Itoa(defval))
	res, err := strconv.Atoi(str)
	if err == nil {
		return res
	}
	return defval
}
func getStringFromDB(sqlx string, defval string) string {

	rows, err := DBH.Query(sqlx)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if rows.Next() {
		var val string
		rows.Scan(&val)
		return val
	}
	return defval
}

func establishDatabase() bool {
	var dbok bool
	var dbi int

	sqlx := "SELECT DBInitialised FROM config"
	rows, err := DBH.Query(sqlx)
	if err == nil && rows.Next() {
		err = rows.Scan(&dbi)
		checkerr(err)
		dbok = dbi > 0
		rows.Close()
		return dbok
	}

	fmt.Println("Warning: establishing basic database")

	_, err = DBH.Exec(chasmsql)
	dbok = err == nil
	if err != nil {
		fmt.Printf("DB establish failed %v\n", err)
	}
	return dbok

}
