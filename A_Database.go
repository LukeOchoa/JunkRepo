package main

import (
	//"github.com/lib/pq"
	"bytes"
	"database/sql"
	"encoding/json"
	"io"

	//"reflect"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const (
	host     = "localhost"
	port     = 5432
	userx    = "luke"
	password = "free144"
	dbname   = "breaker"
)

// Database Misc Functions \\
func create_DB_Connection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s "+
		"sslmode=disable",
		host, port, userx, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// fmt.Println("Successfully connected!")

	return db
}

func getValidIDstr(dbName string) string {
	db := create_DB_Connection()
	// sql string to get all primary key ids
	sql := fmt.Sprintf(`
                SELECT id From %s;`, dbName)

	var idKeys []int
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		idKeys = append(idKeys, id)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
	db.Close()

	var booly bool = false
	var lowestNewValue int = 0
	for _ = range idKeys {
		for y := range idKeys {
			if lowestNewValue == idKeys[y] {
				booly = true
				break
			}
		}
		if booly {
			lowestNewValue = lowestNewValue + 1
		}
		booly = false
	}

	return strconv.Itoa(lowestNewValue)
}

// \\ _________________________________ //

// Supportive Functions \\

// Remove a certain amount of characters from either the beginning or end of a string
func reduceString(str *string, amount int, position string) {
	// TODO This function can still break itself if the user puts in incorrect values
	//      insde the parameters. Use with caution for now

	// use this function to remove from the end or beginning of a string

	var str2 []string
	// Create a slice of the original string
	for _, char := range *str {
		str2 = append(str2, fmt.Sprintf("%c", char))
	}

	// remove part of the string
	if position == "end" {
		str2 = str2[:len(str2)-amount]
	}
	if position == "start" {
		str2 = str2[amount:]
	}

	// convert it back to a real string from a slice
	var realstr string
	for x := range str2 {
		realstr = realstr + str2[x]
	}

	// write the new string to the orignal through a dereferenced pointer
	*str = realstr
}

// removes brackets from []string and returns them in a type string variable
func removeBracketsFromArray(arr []string) string {
	var str string
	for _, v := range arr {
		str = str + v + " "
	}
	reduceString(&str, 1, "end")

	return str
}

func MapString(oldSlicey []string, f func(string) string) []string {
	slicey := make([]string, len(oldSlicey))
	for index, value := range oldSlicey {
		slicey[index] = f(value)
	}
	return slicey
}

// Return a string formated for using in sql strings. With or without quotes ` ' ' `
func sql_Format(str string, quotes bool) string {
	var realstr string
	var quoties string = ""
	if quotes {
		quoties = `'`
	}
	realstr = fmt.Sprintf(`%s%s%s, `, quoties, str, quoties)
	return realstr
}

func hashIt(psswd string) string {

	bs, err := bcrypt.GenerateFromPassword([]byte(psswd), bcrypt.MinCost)
	if err != nil {
		fmt.Println("Internal Server Error! Inside hashIt()...")
		panic(err)
	}
	return string(bs)
}

func (profile nProfile) Marshal() []byte {
	response, err := json.Marshal(profile)
	if err != nil {
		panic(err)
	}
	return response
}

func (profile nProfile) Encode() []byte {
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(&profile)
	return buffer.Bytes()
}

func (profile *nProfile) Decode(reader io.Reader) {
	decoder := json.NewDecoder(reader)
	err2 := decoder.Decode(&profile)
	if err2 != nil {
		panic(err2)
	}
}

func Decode(profile *nProfile, reader io.Reader) {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(profile)
	if err != nil {
		panic(err)
	}
}

// \\ _________________________________ //

func loadUser(sID string, profile *nProfile) bool {
	// Create database connection
	db := create_DB_Connection()

	// Retrive username associated with sessionid(the cookie.Value/sID) from DB
	sql_select := fmt.Sprintf(`
                SELECT username FROM sessions WHERE sessionid='%s';
                `, sID)
	var sessionUsername string
	err := db.QueryRow(sql_select).Scan(&sessionUsername)
	if err != nil {
		fmt.Println("Panic1 inside loadUser()! No username associated with this sessionid: ", sID)
		panic(err)
	}

	// Retrive user's profile from database
	sql_select_2 := fmt.Sprintf(`
                SELECT username, password, firstname, lastname, role FROM userprofile WHERE username='%s';
                `, sessionUsername)
	var username string
	var password string
	var firstname string
	var lastname string
	var role string
	// Scan values into the user's profile struct
	err2 := db.QueryRow(sql_select_2).Scan(
		&username,
		&password,
		&firstname,
		&lastname,
		&role,
	)
	profile.Username = username
	profile.Password = password
	profile.Firstname = firstname
	profile.Lastname = lastname
	profile.Role = role

	if err2 != nil {
		fmt.Println("Panic2 inside loadUser()! Failed to retrive user's profile from database!")
		panic(err2)
	}

	db.Close()

	return true
}

func selectFromDB(column string, table string, condition string, stored_value string) string {

	//SELECT column FROM table WHERE condition ='stored_value';
	db := create_DB_Connection()

	sql_select := fmt.Sprintf(`
                SELECT %s FROM %s WHERE %s='%s'`, column, table, condition, stored_value)
	var database_value string
	err := db.QueryRow(sql_select).Scan(&database_value)

	if err != nil {
		fmt.Println("Panic in selectFromDB()! At sql select statement...")
		panic(err)
	}
	db.Close()

	return database_value
}

func checkIfExists(table string, column string, value string) bool {
	// does (value string) exist?
	db := create_DB_Connection()

	sql_select := fmt.Sprintf(`
			SELECT exists (SELECT 1 FROM %s WHERE %s='%s');`, table, column, value)
	var avail bool
	err := db.QueryRow(sql_select).Scan(&avail)
	if err != nil {
		fmt.Println("Panic in checkIfExists()! At sql select statement...")
		panic(err)
	}
	db.Close()

	// Return true for "it exists"
	return avail
}

func checkUsernameAvailability(profile nProfile) bool {
	db := create_DB_Connection()

	sql_select := fmt.Sprintf(
		`SELECT exists(SELECT 1 FROM userprofile WHERE username='%s');`, profile.Username)
	var avail bool
	err := db.QueryRow(sql_select).Scan(&avail)
	if err != nil {
		panic(err)
	}

	return !avail
}

func dbCreate(crud Crud) {

	drop := removeBracketsFromArray

	// Convert the arrays to formated strings
	table := crud.table
	columns := drop(MapString(crud.column, func(v string) string { return sql_Format(v, false) }))
	columns_values := drop(MapString(crud.column_value, func(v string) string { return sql_Format(v, true) }))
	reduceString(&columns, 2, "end")
	reduceString(&columns_values, 2, "end")

	sql_create := fmt.Sprintf(
		`INSERT INTO %s(%s) VALUES(%s);`, table, columns, columns_values)

	// Insert to database
	db := create_DB_Connection()
	_, err := db.Exec(sql_create)
	if err != nil {
		fmt.Println("Panic inside dbCreate()! At sql insert statement...")
		panic(err)
	}
}

func dbRead(crud Crud) string {

	drop := removeBracketsFromArray

	// Convert the storage variables' values to formated strings
	table := crud.table
	columns := drop(MapString(crud.column, func(v string) string { return sql_Format(v, false) }))
	where := sql_Format(crud.where, false)
	where_condition := sql_Format(crud.where_condition, true)
	reduceString(&columns, 2, "end")
	reduceString(&where, 1, "end")
	reduceString(&where_condition, 1, "end")

	sql_create := fmt.Sprintf(
		`SELECT %s FROM %s`, columns, table)
	sql_part := fmt.Sprintf(
		` WHERE %s=%s`, where, where_condition)
	if crud.where != "" {
		sql_create = sql_create + sql_part + `;`
	}
	fmt.Println(sql_create)
	return sql_create
}

func dbUpdate(crud Crud) {
	// UPDATE table SET column = column_value WHERE condition = where_condition;
	// Convert crud variables to formated strings
	var set string = ""
	for index := 0; index < len(crud.column); index++ {
		set = set + (crud.column[index] + "=" + sql_Format(crud.column_value[index], true))
	}
	reduceString(&set, 2, "end")

	sql_update := fmt.Sprintf(
		`UPDATE %s SET %s WHERE %s='%s';`, crud.table, set, crud.where, crud.where_condition)

	// Insert to database
	db := create_DB_Connection()
	_, err := db.Exec(sql_update)
	if err != nil {
		fmt.Println("Panic inside dbUpdatex()! at UPDATE statement...")
		panic(err)
	}
	db.Close()
}

// TODO make a delete function
func dbDelete(crud Crud) {
	sql_delete := fmt.Sprintf(
		`DELETE FROM %s WHERE %s='%s';`, crud.table, crud.column[0], crud.column_value[0],
	)
	db := create_DB_Connection()
	_, err := db.Exec(sql_delete)
	if err != nil {
		panic(err)
	}
	db.Close()
}
