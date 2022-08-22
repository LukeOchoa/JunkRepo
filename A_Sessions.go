package main

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func initialLoadUser(w http.ResponseWriter, r *http.Request, profile *nProfile) {

	// Retreive User session
	cookie, err := r.Cookie("session")
	if err == nil { // if found
		if checkIfExists("sessions", "sessionid", cookie.Value) { // if found in database
			loadUser(cookie.Value, profile) // assign all data to profile struct
		}
	}
}

func deleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		panic(err)
	}

	crud := Crud{
		table:        "sessions",
		column:       []string{"sessionid"},
		column_value: []string{cookie.Value},
	}
	dbDelete(crud)

	cookie = &http.Cookie{
		Name:     `session`,
		Value:    cookie.Value,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

func createSession(sessionLength int, w http.ResponseWriter, req *http.Request) *http.Cookie {
	// Create uuid for client and assign it to cookie along with max-age and name
	sID := uuid.NewV4()
	cookie := &http.Cookie{
		Name:     `session`,
		Value:    sID.String(),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
	}
	// Write cookie
	http.SetCookie(w, cookie)
	return cookie
}

func testSession2(sessionLengthx int, w http.ResponseWriter, req *http.Request) *http.Cookie {
	sID := "99cf7294-119f-4abc-b4c2-5d84b37bac10"
	cookie := &http.Cookie{
		Name:     `session`,
		Value:    sID,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour),
	}
	http.SetCookie(w, cookie)
	req.AddCookie(cookie)
	return cookie
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	// Retrive the session. If there is no session, return false
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}

	// Search for cookie.Value(sessionid) inside database
	if !checkIfExists("sessions", "sessionid", cookie.Value) {
		// return false for Client is not logged in
		return false
	} else {
		crud := Crud{
			table:           "sessions",
			column:          []string{"lastactivity"},
			column_value:    []string{time.Now().String()},
			where:           "sessionid",
			where_condition: cookie.Value,
		}
		dbUpdate(crud)

		// Refresh the session cookie
		cookie.Expires = time.Now().Add(time.Hour)
		http.SetCookie(w, cookie)

		// return true for Client is already logged in
		return true
	}
}

func cleanSessions() { // comment block properly before you forget what anything does or why...
	for range time.Tick(time.Second * 1) {
		crud := Crud{
			table:  "sessions",
			column: []string{"id", "lastactivity"},
		}
		sql_select := dbRead(crud)

		db := create_DB_Connection()
		var lastactivitys [][]string
		rows, err := db.Query(sql_select)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var lastactivity string
			var id string
			err = rows.Scan(&id, &lastactivity)
			if err != nil {
				fmt.Println("Panic inside cleanSessions() @ rows.Scan()!!! ...")
				panic(err)
			}
			lastactivitys = append(lastactivitys, []string{id, lastactivity})
		}
		err = rows.Err()
		if err != nil {
			fmt.Println("Panic in cleanSessions() @rows.Err()...")
			panic(err)
		}
		rows.Close()
		for _, value := range lastactivitys {
			timey, err2 := time.Parse(RFC3339, value[1])
			if err2 != nil {
				panic(err2)
			}
			if time.Since(timey) > (time.Second * 30) { // does this statement actually work?
				sql_delete := fmt.Sprintf(`DELETE FROM sessions WHERE id=%s;`, value[0])
				_, err3 := db.Exec(sql_delete)
				if err3 != nil {
					fmt.Println("Panic inside cleanSessions() @sql_delete! ...")
					panic(err3)
				}
			}
		}
		db.Close()
	}
}
