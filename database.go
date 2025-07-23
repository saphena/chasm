package main

import "strconv"

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
