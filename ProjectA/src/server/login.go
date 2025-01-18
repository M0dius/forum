package server

import (
	"database/sql"
	"fmt"
	UUID "forum/src/security"
	"log"
	"net/http"
	"time"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println("Redirecting to Home page")
		err := ErrorPageData{Code: "404", ErrorMsg: "PAGE NOT FOUND"}
		errHandler(w, r, &err)
		return
	}

	db, err := sql.Open("sqlite3", "./database/main.db")
	if err != nil {
		log.Println("Database connection failed")
		err := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}
	defer db.Close()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if r.Method == "POST" {
		var userID int
		var dbPassword, userName string
		err = db.QueryRow("SELECT userid, password, username FROM user WHERE email = ?", email).Scan(&userID, &dbPassword, &userName)
		if err != nil {
			if err == sql.ErrNoRows {
				// No credentials found with the given email
				err = templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
					"ErrorMsg": "Invalid email or password",
				})
				if err != nil {
					log.Println("Error rendering login page:", err)
					errData := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
					errHandler(w, r, &errData)
				}
				return
			}
			log.Println("Failed to fetch user data")
			err := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &err)
			return
		}

		// Check if the password is correct
		if !VerifyPassword(password, dbPassword) {
			err := templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
				"ErrorMsg": "Invalid email or password",
			})
			if err != nil {
				log.Println("Error rendering login page:", err)
				errData := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
			}
			return
		}

		// Generate a new session token
		sessionToken, err := UUID.GenerateToken()
		if err != nil {
			log.Println("Error generating UUID:", err)
			errData := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
		}

		stringToken := sessionToken.String()

		//Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    stringToken,
			Expires:  time.Now().Add(1 * time.Hour), //1 hour lifetime
			HttpOnly: true,
		})

		//Set username cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "dotcom_user",
			Value:    userName,
			Expires:  time.Now().Add(1 * time.Hour),
			HttpOnly: true,
		})

		// Update the user's session ID in the session table
		result, err := db.Exec("UPDATE session SET sessionid = ? WHERE userid = ?", stringToken, userID)
		if err != nil {
			log.Println("Error updating session ID in session table:", err)
			errData := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		} else if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected == 0 { //only insert a new row if no record is updated (i.e., no session is found)
			_, err := db.Exec("INSERT INTO session (sessionid, userid, endtime) VALUES (?, ?, ?) RETURNING sessionid",
				stringToken, userID, time.Now().Add(1 * time.Hour))
			if err != nil {
				log.Println("Error creating new session:", err)
				errData := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
				return
			}
		}

		// Update the user's session ID in the database
		_, err = db.Exec("UPDATE user SET current_session = ? WHERE userid = ?", stringToken, userID)
		if err != nil {
			log.Println("Error updating session ID in user table:", err)
			errData := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		}

		log.Println("User logged in with userID:", userID)

		// If login is successful, redirect to the Home page with user ID
		log.Println("Redirecting to Home page with user ID")
		http.Redirect(w, r, fmt.Sprintf("/home?tab=posts&filter=all", userID), http.StatusSeeOther)
	}

	err = templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Println("Error rendering login page:", err)
		errData := ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &errData)
	}
}
